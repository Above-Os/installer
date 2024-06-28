package ctl

import (
	"bytetrade.io/web3os/installer/cmd/ctl/api"
	"bytetrade.io/web3os/installer/cmd/ctl/checksum"
	"bytetrade.io/web3os/installer/cmd/ctl/system"
	"github.com/spf13/cobra"
)

func NewDefaultCommand() *cobra.Command {
	cmds := &cobra.Command{
		Use:               "installer",
		Short:             "Terminus Installer",
		Long:              `......`,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	cmds.AddCommand(api.NewCmdApi())
	cmds.AddCommand(system.NewCmdSystem())
	cmds.AddCommand(checksum.NewCmdChecksum())

	return cmds
}
