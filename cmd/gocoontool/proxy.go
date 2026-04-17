package main

import (
	"context"
	"fmt"

	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	"github.com/tonkeeper/tongo/ton"
)

func cmdProxyState() {
	addr := mustArg(3, "proxy state <address>")
	ctx := context.Background()
	lc := mustLiteClient()

	pc := abiCocoon.NewCocoonProxy(lc, lc).WithAccountId(ton.MustParseAccountID(addr))
	data, err := pc.GetCocoonProxyData(ctx)
	if err != nil {
		fatalf("get proxy data: %v", err)
	}
	fmt.Printf("ProxyState (%s):\n", addr)
	fmt.Printf("  OwnerAddress:      %v\n", data.OwnerAddress)
	fmt.Printf("  ProxyPublicKey:    %v\n", data.ProxyPublicKey)
	fmt.Printf("  RootAddress:       %v\n", data.RootAddress)
	fmt.Printf("  State:             %v\n", data.State)
	fmt.Printf("  Balance:           %v\n", data.Balance)
	fmt.Printf("  Stake:             %v\n", data.Stake)
	fmt.Printf("  UnlockTs:          %v\n", data.UnlockTs)
	fmt.Printf("  PricePerToken:     %v\n", data.PricePerToken)
	fmt.Printf("  WorkerFeePerToken: %v\n", data.WorkerFeePerToken)
	fmt.Printf("  MinProxyStake:     %v\n", data.MinProxyStake)
	fmt.Printf("  MinClientStake:    %v\n", data.MinClientStake)
	fmt.Printf("  ParamsVersion:     %v\n", data.ParamsVersion)
}
