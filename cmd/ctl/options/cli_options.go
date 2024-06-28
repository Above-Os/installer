package options

import "github.com/spf13/cobra"

type CliTerminusUninstallOptions struct {
	Proxy    string
	KubeType string
}

func NewCliTerminusUninstallOptions() *CliTerminusUninstallOptions {
	return &CliTerminusUninstallOptions{}
}

func (o *CliTerminusUninstallOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Proxy, "proxy", "", "Set proxy address, e.g., http://192.168.50.32 or https://your-proxy-domain")
	cmd.Flags().StringVar(&o.KubeType, "type", "k3s", "Specify the container orchestration platform type, please set to k8s or k3s")
}