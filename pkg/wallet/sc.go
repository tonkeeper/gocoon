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
					Seqno:       0,
					SubwalletId: 0,
					PublicKey:   tlb.Uint256(*new(big.Int).SetBytes(w.PublicKey)),
					Status:      0,
					OwnerAddress: tlb.InternalAddress{
						Workchain: int8(w.OwnerAddress.Workchain),
						Address:   w.OwnerAddress.Address,
					},
				},
			},
		},
	}, nil
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
	// Determine current seqno; detect whether wallet SC is already deployed.
	var seqno uint32
	var needInit bool
	if seqnoInt, err := abiCocoon.GetSeqno(ctx, lc, w.Address); err != nil {
		// Wallet SC not yet deployed — use seqno=0 and include StateInit.
		needInit = true
	} else {
		seqno = uint32((*big.Int)(&seqnoInt).Uint64())
	}

	if !code.Exists {
		panic("code must exist")
	}

	//secHash := sha256.New()
	//secHash.Write([]byte(secret))
	//secHashBytes := big.NewInt(0).SetBytes(secHash.Sum(nil))

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
					OwnerAddress:   w.OwnerAddress.ToInternal(),
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
		SendExcessesTo: w.OwnerAddress.ToInternal(),
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

	// Marshal the internal message to a cell (ForwardMsg.Msg is a boc.Cell).
	msgCell := boc.NewCell()
	if err := tlb.Marshal(msgCell, registerMsg); err != nil {
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

	// Sign: hash of the serialized ExternalSignedMessage.
	signCell := boc.NewCell()
	if err := tlb.Marshal(signCell, signedMsg); err != nil {
		return fmt.Errorf("marshal for signing: %w", err)
	}
	extMsg, err := signedMsg.Sign(w.PrivateKey)
	if err != nil {
		return fmt.Errorf("sign external message: %w", err)
	}

	//celll, _ := extMsg.ToCell()
	//xxxx, _ := celll.ToBocBase64()
	//fmt.Println(xxxx)

	// Include StateInit when wallet SC is not yet deployed.
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

	// Serialize message to BOC.
	cell := boc.NewCell()
	if err := tlb.Marshal(cell, tlbMsg); err != nil {
		return fmt.Errorf("marshal message to cell: %w", err)
	}
	msgBoc, err := cell.ToBoc()
	if err != nil {
		return fmt.Errorf("serialize to BOC: %w", err)
	}

	refCell := boc.MustDeserializeSinglRootHex("b5ee9c7201021f010005250003b57fcb6f4b0298382d13422f87dad82bcbf26541b90f97f1004589c0e1e8d981a4e00003f00d723d681f872eb5eb1ff0bc387691d6281612d298266d069d414d52357f3d4ab8054887b00003f009ea2819969d6bef6000246a47a4880102030201e004050082723d9d8deb512fe003f6229eeba3794626f4e0fae797453cb497e599db02807e401a3140fb656699b64c9f4b14da8f48208e9c068ed478f9398d550e1ca65f13b2020f0c420618a83684401d1e02458801f96de960530705a26845f0fb5b05797e4ca83721f2fe2008b1381c3d1b30349c1e06070101df1b0201340809019a8d2cdeb4864c20c45b2e5f6d31f0ccea6ef8cf740d6d5a5a5634767fc23465f88403804a0986ea422a53ec63319c8eca3968ea0230ab4242ea0d27c705516b0a0000000069d6cd0600000000001a0114ff00f4a413f4bcf2c80b0a009b0000000000000000421589ec5082a8ed36e08912517bcc8f2eb286bdff059896ae040c8b0a2560e6000000008011c3299746b0061a3caa2ee199d0c27059446e3ddbbb67ad563867864df9d99a700201200b0c0201480d0e00e0f28308d71820d31fd31fd31f02f823bbf2d406ed44d0d31fd31fd3ffd31ffa40d12171b0f2d4075154baf2e4085162baf2e40906f901541076f910f2e40af8276f2230821077359400b9f2d40bf800029320d74a96d307d402fb00e83001a4c8cb1f14cb1f12cbffcb1f01cf16c9ed540202ca0f10020120161702012011120049a9f21402b3c5940233c585b2fff2413232c05400fe80807e80b2cfc4b2c7c4b2fff333326001f5d3b68bb7edb088831c02456f8007434c0c05c6c2456f83e900c0074c7c86084095964d32e88a08431669f34eeac48a084268491f02eac6497c0f83b513434c7f4c7f4fff4c7fe903454dc31c17cb90409a084271a7cddaea78415d7c1f4cfcc74c1f50c007ec03801b0003cb9044134c1f448301dc8701880b01d613020276141500ea5312b121b1f2e411018e295f07820898968072fb0280777080185410337003c8cb0558cf1601fa02cb6a12cb1fcb07c98306fb00e0378e19350271b101c8cb1f12cb1f13cbff12cb1f01cf16c9ed54db31e0058e1d028210fffffffeb001c8cb1f12cb1f13cbff12cb1f01cf16c9ed54db31e05f0500691cf232c1c044440072c7c7b2c7c732c01402be8094023e8085b2c7c532c7c4b2c7f2c7f2c7f2c7c07e80807e80bd003d003d00326000553434c1c07000fcb8fc34c7f4c7f4c03e803e8034c7f4c7f4c7f4c7f4c7f4c7fe803e803d013d013d010c200201201819001bbdfddf6a2684080b06b90fd201840017bb39ced44d0d33f31d70bff80011b8c97ed44d0d70b1f8016862004e1152573f579bb55f437b26d365e3b9a8e5045e0cb4cfc568a33c206e30a180a1dcd65000000000000000000000000000011c01b16801f96de960530705a26845f0fb5b05797e4ca83721f2fe2008b1381c3d1b30349d002708a92b9fabcddaafa1bd9369b2f1dcd472822f065a67e2b4519e10371850c050ee45e040060c6d4800007e01ae47ad04d3ad7decc01c006bc45f9f3b37f3c00d5a22c9e56a868b186e7d755b801f96de960530705a26845f0fb5b05797e4ca83721f2fe2008b1381c3d1b30349d0009d41aea313880000000000000000124000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020006fc984a8f84c0c6d3000000000000200000000000396e2d2190cbde207d311786286ae9abbae7ba60094089dc072b993742c30534440902364")

	fmt.Println(refCell)

	bocB64 := base64.StdEncoding.EncodeToString(msgBoc)
	if err := sender.SendBlockchainMessage(ctx, &tonapi.SendBlockchainMessageReq{
		Boc: tonapi.NewOptString(bocB64),
	}); err != nil {
		return fmt.Errorf("send blockchain message: %w", err)
	}
	logger.Info("register tx sent")
	return nil
}
