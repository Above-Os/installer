package options

import "github.com/spf13/cobra"

type ApiOptions struct {
	Enabled  bool
	Port     string
	LogLevel string
	Proxy    string
}

func NewApiOptions() *ApiOptions {
	return &ApiOptions{}
}

func (o *ApiOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.Enabled, "enabled", false, "running api server")
	cmd.Flags().StringVar(&o.Port, "port", ":30180", "listen port")
	cmd.Flags().StringVar(&o.LogLevel, "log-level", "debug", "log level")
	cmd.Flags().StringVar(&o.Proxy, "proxy", "", "proxy")
}
