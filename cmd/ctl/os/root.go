package os

import "github.com/spf13/cobra"

func NewCmdOs() *cobra.Command {
	rootOsCmd := &cobra.Command{
		Use:   "terminus",
		Short: "Terminus install, uninstall or restore",
	}

	rootOsCmd.AddCommand(NewCmdInstallOs())
	rootOsCmd.AddCommand(NewCmdRestoreOs())
	rootOsCmd.AddCommand(NewCmdUninstallOs())
	return rootOsCmd
}
