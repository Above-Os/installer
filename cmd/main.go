package main

import (
	"fmt"
	"os"

	"bytetrade.io/web3os/installer/cmd/ctl"
)

func main() {
	cmd := ctl.NewDefaultCommand()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
