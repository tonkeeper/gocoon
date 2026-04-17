package main

import (
	"context"
	"fmt"

	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	"github.com/tonkeeper/tongo/ton"
)

func cmdClientState() {
	addr := mustArg(3, "client state <address>")
	ctx := context.Background()
	lc := mustLiteClient()

	cc := abiCocoon.NewCocoonClient(lc, lc).WithAccountId(ton.MustParseAccountID(addr))
	data, err := cc.GetCocoonClientData(ctx)
	if err != nil {
		fatalf("get client data: %v", err)
	}
	fmt.Printf("ClientState (%s):\n", addr)
	fmt.Printf("  OwnerAddress:   %v\n", data.OwnerAddress)
	fmt.Printf("  ProxyAddress:   %v\n", data.ProxyAddress)
	fmt.Printf("  ProxyPublicKey: %v\n", data.ProxyPublicKey)
	fmt.Printf("  State:          %v\n", data.State)
	fmt.Printf("  Balance:        %v\n", data.Balance)
	fmt.Printf("  Stake:          %v\n", data.Stake)
	fmt.Printf("  TokensUsed:     %v\n", data.TokensUsed)
	fmt.Printf("  UnlockTs:       %v\n", data.UnlockTs)
	fmt.Printf("  SecretHash:     %v\n", data.SecretHash)
}
