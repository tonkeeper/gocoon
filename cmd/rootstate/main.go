package main

import (
	"context"
	"fmt"
	"log"

	tongo "github.com/tonkeeper/tongo"
	abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"

	//abiCocoon "github.com/tonkeeper/tongo/abi-tolk/abiGenerated/cocoon"
	"github.com/tonkeeper/tongo/liteapi"
)

func main() {
	const rootContractAddr = "EQCns7bYSp0igFvS1wpb5wsZjCKCV19MD5AVzI4EyxsnU73k"

	client, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		log.Fatalf("create liteapi client: %v", err)
	}
	rootAddr := tongo.MustParseAddress(rootContractAddr)

	rootClient := abiCocoon.NewCocoonRoot(client, client).WithAccountId(rootAddr.ID)

	rootClient.GetLastProxySeqno(context.TODO())

	ctx := context.Background()

	rootStore, err := rootClient.Storage(ctx)
	if err != nil {
		log.Fatalf("get root storage: %v", err)
	}
	fmt.Printf("RootStorage: %v\n", rootStore.Data.Value.RegisteredProxies.Values()[0].Address)
	if true {
		return
	}

	rootData, err := rootClient.GetCocoonData(ctx)
	if err != nil {
		log.Fatalf("get cocoon data: %v", err)
	}
	fmt.Printf("CocoonData:\n")
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

	lastProxySeqno, err := rootClient.GetLastProxySeqno(ctx)
	if err != nil {
		log.Fatalf("get last proxy seqno: %v", err)
	}
	fmt.Printf("\nLastProxySeqno: %v\n", lastProxySeqno)

	curParams, err := rootClient.GetCurParams(ctx)
	if err != nil {
		log.Fatalf("get cur params: %v", err)
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

	workerAddr := tongo.MustParseAddress("UQBksjMUcEwjSXsi-_XRaQAJJd7NsCCEUNp0kgAPk0Tr-IrC").ID
	workerClient := abiCocoon.NewCocoonWorker(client).WithAccountId(workerAddr)
	workerData, err := workerClient.GetCocoonWorkerData(ctx)
	if err != nil {
		log.Fatalf("get cocoon worker data: %v", err)
	}
	fmt.Printf("\nCocoonWorkerData (%v):\n", workerAddr)
	fmt.Printf("  OwnerAddress:   %v\n", workerData.OwnerAddress)
	fmt.Printf("  ProxyAddress:   %v\n", workerData.ProxyAddress)
	fmt.Printf("  ProxyPublicKey: %v\n", workerData.ProxyPublicKey)
	fmt.Printf("  State:          %v\n", workerData.State)
	fmt.Printf("  Tokens:         %v\n", workerData.Tokens)

	proxyAddrs := []tongo.AccountID{
		tongo.MustParseAddress("UQCUKdaSCNCRo3lVTB_6etw5vF6rJWHS9wE5i__dzIhY3-yK").ID,
		tongo.MustParseAddress("UQCfAloyPc1B_VctBTmKsh5fHFL4FqnqMxuuOjf5ornEUUu8").ID,
		tongo.MustParseAddress("UQCH-ucMN3d4XJ8vBe_8OKbx6ZZpfFXESumWl1FedouIlC3Y").ID,
	}

	for _, proxyAddr := range proxyAddrs {
		proxyClient := abiCocoon.NewCocoonProxy(client).WithAccountId(proxyAddr)
		proxyData, err := proxyClient.GetCocoonProxyData(ctx)
		if err != nil {
			log.Fatalf("get cocoon proxy data for %v: %v", proxyAddr, err)
		}
		fmt.Printf("\nCocoonProxyData (%v):\n", proxyAddr)
		fmt.Printf("  OwnerAddress:      %v\n", proxyData.OwnerAddress)
		fmt.Printf("  ProxyPublicKey:    %v\n", proxyData.ProxyPublicKey)
		fmt.Printf("  RootAddress:       %v\n", proxyData.RootAddress)
		fmt.Printf("  State:             %v\n", proxyData.State)
		fmt.Printf("  Balance:           %v\n", proxyData.Balance)
		fmt.Printf("  Stake:             %v\n", proxyData.Stake)
		fmt.Printf("  UnlockTs:          %v\n", proxyData.UnlockTs)
		fmt.Printf("  PricePerToken:     %v\n", proxyData.PricePerToken)
		fmt.Printf("  WorkerFeePerToken: %v\n", proxyData.WorkerFeePerToken)
		fmt.Printf("  MinProxyStake:     %v\n", proxyData.MinProxyStake)
		fmt.Printf("  MinClientStake:    %v\n", proxyData.MinClientStake)
		fmt.Printf("  ParamsVersion:     %v\n", proxyData.ParamsVersion)
	}

	clientClient := abiCocoon.NewCocoonWallet(client).WithAccountId(
		tongo.MustParseAccountID("UQD8tvSwKYOC0TQi-H2tgry_JlQbkPl_EARYnA4ejZgaTqI9"),
	)
	x, _ := clientClient.GetOwnerAddress(ctx)
	fmt.Printf("\nCocoonClientData (%v):\n", x)
}
