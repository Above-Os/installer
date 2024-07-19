package options

import "github.com/spf13/cobra"

// ~ CliTerminusUninstallOptions
type CliTerminusUninstallOptions struct {
	Proxy string
}

func NewCliTerminusUninstallOptions() *CliTerminusUninstallOptions {
	return &CliTerminusUninstallOptions{}
}

func (o *CliTerminusUninstallOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Proxy, "proxy", "", "Set proxy address, e.g., 192.168.50.32 or your-proxy-domain")
}

// ~ CliTerminusInstallOptions
type CliTerminusInstallOptions struct {
	KubeType string
	Proxy    string
}

func NewCliTerminusInstallOptions() *CliTerminusInstallOptions {
	return &CliTerminusInstallOptions{}
}

func (o *CliTerminusInstallOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.KubeType, "kube", "", "Set kube type, e.g., k3s or k8s")
	cmd.Flags().StringVar(&o.Proxy, "proxy", "", "Set proxy address, e.g., 192.168.50.32 or your-proxy-domain")
}
