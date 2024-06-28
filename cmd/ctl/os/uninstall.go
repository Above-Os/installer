package os

import (
	"fmt"

	"bytetrade.io/web3os/installer/cmd/ctl/options"
	"github.com/spf13/cobra"
)

type UninstallOsOptions struct {
	UninstallOptions *options.CliTerminusUninstallOptions
}

func NewUninstallOsOptions() *UninstallOsOptions {
	return &UninstallOsOptions{
		UninstallOptions: options.NewCliTerminusUninstallOptions(),
	}
}

func NewCmdUninstallOs() *cobra.Command {
	o := NewUninstallOsOptions()
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall Terminus",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Uninstall Terminus")
			uninstall()
		},
	}
	o.UninstallOptions.AddFlags(cmd)
	return cmd
}

func uninstall() {

}
