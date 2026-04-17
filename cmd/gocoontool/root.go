package main

import (
	"context"
	"fmt"
	"os"

	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	"github.com/tonkeeper/tongo/ton"
)

const defaultRootAddr = "EQCns7bYSp0igFvS1wpb5wsZjCKCV19MD5AVzI4EyxsnU73k"

func cmdRootState() {
	rootAddrStr := defaultRootAddr
	if len(os.Args) >= 4 {
		rootAddrStr = os.Args[3]
	}

	ctx := context.Background()
	lc := mustLiteClient()
	rootClient := abiCocoon.NewCocoonRoot(lc, lc).WithAccountId(ton.MustParseAccountID(rootAddrStr))

	_, rootStore, err := rootClient.AccountState(ctx)
	if err != nil {
		fatalf("get root storage: %v", err)
	}
	proxies := rootStore.Data.Value.RegisteredProxies.Values()
	fmt.Printf("RootStorage:\n")
	fmt.Printf("  OwnerAddress: %v\n", rootStore.OwnerAddress)
	fmt.Printf("  Version:      %v\n", rootStore.Version)
	fmt.Printf("  Proxies (%d):\n", len(proxies))
	for i, p := range proxies {
		fmt.Printf("    [%d] %s\n", i, p.Address)
	}

	rootData, err := rootClient.GetCocoonData(ctx)
	if err != nil {
		fatalf("get cocoon data: %v", err)
	}
	fmt.Printf("\nCocoonData:\n")
	fmt.Printf("  Version:           %v\n", rootData.Version)
	fmt.Printf("  LastProxySeqno:    %v\n", rootData.LastProxySeqno)
	fmt.Printf("  ParamsVersion:     %v\n", rootData.ParamsVersion)
	fmt.Printf("  UniqueId:          %v\n", rootData.UniqueId)
	fmt.Printf("  IsTest:            %v\n", rootData.IsTest)
	fmt.Printf("  PricePerToken:     %v\n", rootData.PricePerToken)
	fmt.Printf("  WorkerFeePerToken: %v\n", rootData.WorkerFeePerToken)
	fmt.Printf("  MinProxyStake:     %v\n", rootData.MinProxyStake)
	fmt.Printf("  MinClientStake:    %v\n", rootData.MinClientStake)
	fmt.Printf("  OwnerAddress:      %v\n", rootData.OwnerAddress)

	curParams, err := rootClient.GetCurParams(ctx)
	if err != nil {
		fatalf("get cur params: %v", err)
	}
	fmt.Printf("\nCurrentCocoonParams:\n")
	fmt.Printf("  ParamsVersion:                  %v\n", curParams.ParamsVersion)
	fmt.Printf("  UniqueId:                       %v\n", curParams.UniqueId)
	fmt.Printf("  IsTest:                         %v\n", curParams.IsTest)
	fmt.Printf("  PricePerToken:                  %v\n", curParams.PricePerToken)
	fmt.Printf("  WorkerFeePerToken:              %v\n", curParams.WorkerFeePerToken)
	fmt.Printf("  CachedTokensPriceMultiplier:    %v\n", curParams.CachedTokensPriceMultiplier)
	fmt.Printf("  ReasoningTokensPriceMultiplier: %v\n", curParams.ReasoningTokensPriceMultiplier)
	fmt.Printf("  ProxyDelayBeforeClose:          %v\n", curParams.ProxyDelayBeforeClose)
	fmt.Printf("  ClientDelayBeforeClose:         %v\n", curParams.ClientDelayBeforeClose)
	fmt.Printf("  MinProxyStake:                  %v\n", curParams.MinProxyStake)
	fmt.Printf("  MinClientStake:                 %v\n", curParams.MinClientStake)
	fmt.Printf("  ProxyCodeHash:                  %v\n", curParams.ProxyCodeHash)
	fmt.Printf("  WorkerCodeHash:                 %v\n", curParams.WorkerCodeHash)
	fmt.Printf("  ClientCodeHash:                 %v\n", curParams.ClientCodeHash)
}
