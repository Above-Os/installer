package system

import "github.com/spf13/cobra"

func NewCmdSystem() *cobra.Command {
	rootSystemCmd := &cobra.Command{
		Use:   "terminus",
		Short: "Terminus install, uninstall or restore",
	}

	rootSystemCmd.AddCommand(NewCmdSystemCreate())
	rootSystemCmd.AddCommand(NewCmdSystemRestore())
	rootSystemCmd.AddCommand(NewCmdSystemDelete())
	return rootSystemCmd
}
