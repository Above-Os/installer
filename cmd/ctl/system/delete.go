package system

import (
	"fmt"

	"bytetrade.io/web3os/installer/cmd/ctl/options"
	"github.com/spf13/cobra"
)

type SystemDeleteOptions struct {
	DeleteOptions *options.CliTerminusUninstallOptions
}

func NewSystemDeleteOptions() *SystemDeleteOptions {
	return &SystemDeleteOptions{
		DeleteOptions: options.NewCliTerminusUninstallOptions(),
	}
}

func NewCmdSystemDelete() *cobra.Command {
	o := NewSystemDeleteOptions()
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall Terminus",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Uninstall Terminus")
			uninstall()
		},
	}
	o.DeleteOptions.AddFlags(cmd)
	return cmd
}

func uninstall() {

}
