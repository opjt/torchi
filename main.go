package main

import (
	"fmt"
	"os"
	"torchi/cmd"
	"torchi/internal/core"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("use: torchi <command> [options]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "serve":
		core.RunServer(core.Modules)
	case "genkey":
		vKey, err := cmd.GenKey()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("publicKey:\n", vKey.PublicKey)
		fmt.Println("privateKey:\n", vKey.PrivateKey)

	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}
