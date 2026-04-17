package main

import (
	"fmt"
	"os"

	"github.com/tonkeeper/tongo/liteapi"
)

func main() {
	if len(os.Args) < 3 {
		usage()
	}
	cmd := os.Args[1] + " " + os.Args[2]
	switch cmd {
	case "wallet generate":
		cmdWalletGenerate()
	case "wallet deploy":
		cmdWalletDeploy()
	case "wallet state":
		cmdWalletState()
	case "client state":
		cmdClientState()
	case "proxy state":
		cmdProxyState()
	case "root state":
		cmdRootState()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: gocoontool <command>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  wallet generate          generate a new Ed25519 private key")
	fmt.Fprintln(os.Stderr, "  wallet deploy            deploy the cocoon wallet SC on TON")
	fmt.Fprintln(os.Stderr, "                             env: PRIVATE_KEY, OWNER_ADDRESS")
	fmt.Fprintln(os.Stderr, "  wallet state <address>   print wallet SC state")
	fmt.Fprintln(os.Stderr, "  client state <address>   print client SC state")
	fmt.Fprintln(os.Stderr, "  proxy state <address>    print proxy SC state")
	fmt.Fprintln(os.Stderr, "  root state [address]     print root SC state (default: mainnet)")
	os.Exit(1)
}

func mustLiteClient() *liteapi.Client {
	lc, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		fatalf("create liteapi client: %v", err)
	}
	return lc
}

func mustArg(i int, hint string) string {
	if len(os.Args) <= i {
		fmt.Fprintf(os.Stderr, "usage: gocoontool %s\n", hint)
		os.Exit(1)
	}
	return os.Args[i]
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
