package cocoonWallet

import (
	"context"
	"fmt"
	"math/big"
	"time"

	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

// Code is the hardcoded cocoon_wallet contract bytecode.
var Code = boc.MustDeserializeSinglRootHex(`b5ee9c724102110100024b000114ff00f4a413f4bcf2c80b010201200210020148030b0202ca040a020120050701f5d3b68bb7edb088831c02456f8007434c0c05c6c2456f83e900c0074c7c86084095964d32e88a08431669f34eeac48a084268491f02eac6497c0f83b513434c7f4c7f4fff4c7fe903454dc31c17cb90409a084271a7cddaea78415d7c1f4cfcc74c1f50c007ec03801b0003cb9044134c1f448301dc8701880b01d60600ea5312b121b1f2e411018e295f07820898968072fb0280777080185410337003c8cb0558cf1601fa02cb6a12cb1fcb07c98306fb00e0378e19350271b101c8cb1f12cb1f13cbff12cb1f01cf16c9ed54db31e0058e1d028210fffffffeb001c8cb1f12cb1f13cbff12cb1f01cf16c9ed54db31e05f05020276080900691cf232c1c044440072c7c7b2c7c732c01402be8094023e8085b2c7c532c7c4b2c7f2c7f2c7f2c7c07e80807e80bd003d003d00326000553434c1c07000fcb8fc34c7f4c7f4c03e803e8034c7f4c7f4c7f4c7f4c7f4c7fe803e803d013d013d010c200049a9f21402b3c5940233c585b2fff2413232c05400fe80807e80b2cfc4b2c7c4b2fff33332600201200c0f0201200d0e0017bb39ced44d0d33f31d70bff80011b8c97ed44d0d70b1f8001bbdfddf6a2684080b06b90fd2018400e0f28308d71820d31fd31fd31f02f823bbf2d406ed44d0d31fd31fd3ffd31ffa40d12171b0f2d4075154baf2e4085162baf2e40906f901541076f910f2e40af8276f2230821077359400b9f2d40bf800029320d74a96d307d402fb00e83001a4c8cb1f14cb1f12cbffcb1f01cf16c9ed545d2b2126`)

func (w *Wallet) buildStateInit() (*tlb.StateInitT[*abiCocoon.WalletStorage], error) {
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
					PublicKey:    tlb.Uint256(*new(big.Int).SetBytes(w.publicKey)),
					Status:       0,
					OwnerAddress: w.ownerAddress.ToInternal(),
				},
			},
		},
	}, nil
}

// Deploy sends the initial external message with StateInit to deploy the wallet SC on TON.
// Must be called once before ForwardMessage can be used.
func (w *Wallet) Deploy(ctx context.Context) error {
	si, err := w.buildStateInit()
	if err != nil {
		return fmt.Errorf("build state init: %w", err)
	}
	signedMsg := abiCocoon.ExternalSignedMessage{
		SubwalletId: 0,
		ValidUntil:  tlb.Uint32(uint32(time.Now().Unix() + 3600)),
		MsgSeqno:    0,
		Forward:     abiCocoon.ForwardMsgs{},
	}
	return w.broadcast(ctx, signedMsg, si)
}

// ForwardMessage wraps msg in a signed external message and broadcasts it.
// The wallet SC must already be deployed; returns an error if the seqno cannot be fetched.
func (w *Wallet) ForwardMessage(ctx context.Context, msg tlb.Message) error {
	seqnoInt, err := abiCocoon.GetSeqno(ctx, w.lc, w.address)
	if err != nil {
		return fmt.Errorf("get seqno: %w", err)
	}
	seqno := uint32((*big.Int)(&seqnoInt).Uint64())

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
	return w.broadcast(ctx, signedMsg, nil)
}

func (w *Wallet) broadcast(ctx context.Context, signedMsg abiCocoon.ExternalSignedMessage, init *tlb.StateInitT[*abiCocoon.WalletStorage]) error {
	extMsg, err := signedMsg.Sign(w.privateKey)
	if err != nil {
		return fmt.Errorf("sign external message: %w", err)
	}
	tlbMsg, err := extMsg.ToExternal(w.address, init)
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
	_, err = w.lc.SendMessage(ctx, msgBoc)
	return err
}
