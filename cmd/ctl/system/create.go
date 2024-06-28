package system

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdSystemCreate() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Install Terminus",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Create Terminus")
		},
	}
}
