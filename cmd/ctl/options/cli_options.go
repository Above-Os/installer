package options

import "github.com/spf13/cobra"

type CliTerminusUninstallOptions struct {
	Proxy string
}

func NewCliTerminusUninstallOptions() *CliTerminusUninstallOptions {
	return &CliTerminusUninstallOptions{}
}

func (o *CliTerminusUninstallOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Proxy, "proxy", "", "Set proxy address, e.g., 192.168.50.32 or your-proxy-domain")
}
