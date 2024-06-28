package options

import "github.com/spf13/cobra"

type CliTerminusUninstallOptions struct {
	Proxy    string
	KubeType string
	Do       string
}

func NewCliTerminusUninstallOptions() *CliTerminusUninstallOptions {
	return &CliTerminusUninstallOptions{}
}

func (o *CliTerminusUninstallOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Proxy, "proxy", "", "Set proxy address, e.g., 192.168.50.32 or your-proxy-domain")
	cmd.Flags().StringVar(&o.KubeType, "type", "k3s", "Specify the container orchestration platform type, please set to k8s or k3s")
	cmd.Flags().StringVar(&o.Do, "do", "", "Run uninstall")
}
