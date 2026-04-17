package wallet

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	tonapi "github.com/tonkeeper/tonapi-go"
	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tl"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"go.uber.org/zap"
)

// BlockchainSender can broadcast a signed BOC to the TON network.
type BlockchainSender interface {
	SendBlockchainMessage(ctx context.Context, request *tonapi.SendBlockchainMessageReq) error
}

// Code is the hardcoded cocoon_wallet contract bytecode.
var Code = boc.MustDeserializeSinglRootHex(`b5ee9c724102110100024b000114ff00f4a413f4bcf2c80b010201200210020148030b0202ca040a020120050701f5d3b68bb7edb088831c02456f8007434c0c05c6c2456f83e900c0074c7c86084095964d32e88a08431669f34eeac48a084268491f02eac6497c0f83b513434c7f4c7f4fff4c7fe903454dc31c17cb90409a084271a7cddaea78415d7c1f4cfcc74c1f50c007ec03801b0003cb9044134c1f448301dc8701880b01d60600ea5312b121b1f2e411018e295f07820898968072fb0280777080185410337003c8cb0558cf1601fa02cb6a12cb1fcb07c98306fb00e0378e19350271b101c8cb1f12cb1f13cbff12cb1f01cf16c9ed54db31e0058e1d028210fffffffeb001c8cb1f12cb1f13cbff12cb1f01cf16c9ed54db31e05f05020276080900691cf232c1c044440072c7c7b2c7c732c01402be8094023e8085b2c7c532c7c4b2c7f2c7f2c7f2c7c07e80807e80bd003d003d00326000553434c1c07000fcb8fc34c7f4c7f4c03e803e8034c7f4c7f4c7f4c7f4c7f4c7fe803e803d013d013d010c200049a9f21402b3c5940233c585b2fff2413232c05400fe80807e80b2cfc4b2c7c4b2fff33332600201200c0f0201200d0e0017bb39ced44d0d33f31d70bff80011b8c97ed44d0d70b1f8001bbdfddf6a2684080b06b90fd2018400e0f28308d71820d31fd31fd31f02f823bbf2d406ed44d0d31fd31fd3ffd31ffa40d12171b0f2d4075154baf2e4085162baf2e40906f901541076f910f2e40af8276f2230821077359400b9f2d40bf800029320d74a96d307d402fb00e83001a4c8cb1f14cb1f12cbffcb1f01cf16c9ed545d2b2126`)

// BuildStateInit constructs the StateInit for the cocoon_wallet SC using the
// hardcoded contract code and the wallet's public key / owner address.
// The hash of the returned cell is the wallet SC address on workchain 0.
func BuildStateInit(w *Wallet) (*tlb.StateInitT[*abiCocoon.WalletStorage], error) {
	return &tlb.StateInitT[*abiCocoon.WalletStorage]{
		Code: tlb.Maybe[tlb.Ref[boc.Cell]]{
			Exists: true,
			Value:  tlb.Ref[boc.Cell]{Value: *Code},
		},
		Data: tlb.Maybe[tlb.Ref[*abiCocoon.WalletStorage]]{
			Exists: true,
			Value: tlb.Ref[*abiCocoon.WalletStorage]{
				Value: &abiCocoon.WalletStorage{
					Seqno:        0,
					SubwalletId:  0,
					PublicKey:    tlb.Uint256(*new(big.Int).SetBytes(w.PublicKey)),
					Status:       0,
					OwnerAddress: w.OwnerAddress.ToInternal(),
				},
			},
		},
	}, nil
}

// ForwardMessage wraps msg in a signed external message, determines the current
// seqno (including StateInit when the wallet SC is not yet deployed), and
// broadcasts the resulting BOC via sender.
func (w *Wallet) ForwardMessage(ctx context.Context, lc *liteapi.Client, msg tlb.Message, sender BlockchainSender) error {
	var seqno uint32
	var needInit bool
	if seqnoInt, err := abiCocoon.GetSeqno(ctx, lc, w.Address); err != nil {
		needInit = true
	} else {
		seqno = uint32((*big.Int)(&seqnoInt).Uint64())
	}

	msgCell := boc.NewCell()
	if err := tlb.Marshal(msgCell, msg); err != nil {
		return fmt.Errorf("marshal internal msg: %w", err)
	}

	signedMsg := abiCocoon.ExternalSignedMessage{
		SubwalletId: 0,
		ValidUntil:  tlb.Uint32(uint32(time.Now().Unix() + 3600)),
		MsgSeqno:    tlb.Uint32(seqno),
		Forward: abiCocoon.ForwardMsgs{{
			Mode: 0,
			Msg:  *msgCell,
		}},
	}

	extMsg, err := signedMsg.Sign(w.PrivateKey)
	if err != nil {
		return fmt.Errorf("sign external message: %w", err)
	}

	var init *tlb.StateInitT[*abiCocoon.WalletStorage]
	if needInit {
		si, err := BuildStateInit(w)
		if err != nil {
			return fmt.Errorf("build wallet state init: %w", err)
		}
		init = si
	}

	tlbMsg, err := extMsg.ToExternal(w.Address, init)
	if err != nil {
		return fmt.Errorf("build external message: %w", err)
	}

	cell := boc.NewCell()
	if err := tlb.Marshal(cell, tlbMsg); err != nil {
		return fmt.Errorf("marshal message to cell: %w", err)
	}
	msgBoc, err := cell.ToBoc()
	if err != nil {
		return fmt.Errorf("serialize to BOC: %w", err)
	}

	bocB64 := base64.StdEncoding.EncodeToString(msgBoc)
	return sender.SendBlockchainMessage(ctx, &tonapi.SendBlockchainMessageReq{
		Boc: tonapi.NewOptString(bocB64),
	})
}

// SendRegisterTx builds and broadcasts the OwnerClientRegister internal message
// from the cocoon wallet SC to the client SC, triggering the on-chain registration
// that long auth waits for.
func SendRegisterTx(
	ctx context.Context, lc *liteapi.Client, w *Wallet,
	clientScAddr ton.AccountID, nonce uint64, sender BlockchainSender,
	code tlb.Maybe[boc.Cell], secret string, logger *zap.Logger,
	proxyAddr ton.AccountID, proxyPubKey tl.Int256, minclstake tlb.Coins,
	cocoonParams *abiCocoon.CocoonParams) error {
	_ = secret // reserved for future use

	if !code.Exists {
		panic("code must exist")
	}

	x := big.NewInt(0).SetBytes(proxyPubKey[:])
	proxyPubKey1 := tlb.Uint256(*x)

	zero := big.NewInt(0)
	cocoonParamsCopy := *cocoonParams
	cocoonParamsCopy.ClientScCode.Exists = false
	cocoonParamsCopy.WorkerScCode.Exists = false
	cocoonParamsCopy.ProxyScCode.Exists = false

	clientInitState := tlb.StateInitT[*abiCocoon.ClientStorage]{
		Code: tlb.JustRef(code.Value),
		Data: tlb.JustRef(&abiCocoon.ClientStorage{
			State:      0,
			Balance:    0,
			Stake:      minclstake,
			TokensUsed: 0,
			UnlockTs:   0,
			SecretHash: tlb.Uint256(*zero),
			ConstDataRef: tlb.RefT[*abiCocoon.ClientConstData]{
				Value: &abiCocoon.ClientConstData{
					OwnerAddress:   w.Address.ToInternal(),
					ProxyAddress:   proxyAddr.ToInternal(),
					ProxyPublicKey: tlb.Uint256(proxyPubKey1),
				},
			},
			Params: tlb.RefT[*abiCocoon.CocoonParams]{Value: &cocoonParamsCopy},
		}),
	}

	cliInStCell, _ := clientInitState.ToCell()
	cliInStCellHash, _ := cliInStCell.HashString()
	fmt.Println("cliInStCellHash", cliInStCellHash)

	// Build OwnerClientRegister internal message to the client SC.
	registerMsg, err := abiCocoon.OwnerClientRegister{
		QueryId:        0,
		Nonce:          tlb.Uint64(nonce),
		SendExcessesTo: w.Address.ToInternal(),
	}.ToInternal(
		tlb.InternalAddress{
			Workchain: int8(clientScAddr.Workchain),
			Address:   tlb.Bits256(clientScAddr.Address),
		},
		tlb.Grams(50_000_000)+minclstake, // 0.05 TON — covers gas + a forwarded message to proxy SC
		true,
		&clientInitState,
	)
	if err != nil {
		return fmt.Errorf("build OwnerClientRegister: %w", err)
	}

	if err := w.ForwardMessage(ctx, lc, registerMsg, sender); err != nil {
		return fmt.Errorf("forward register message: %w", err)
	}
	logger.Info("register tx sent")
	return nil
}
