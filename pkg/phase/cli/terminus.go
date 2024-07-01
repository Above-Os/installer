package cli

import (
	"bytetrade.io/web3os/installer/cmd/ctl/options"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/cache"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/pipeline"
	"bytetrade.io/web3os/installer/pkg/terminus"
)

func Uninstall(o *options.CliTerminusUninstallOptions) error {
	runtime, err := common.NewLocalRuntime(false, false)
	if err != nil {
		return err
	}

	m := []module.Module{
		&terminus.UninstallTerminusCliModule{},
		&terminus.ExecUninstallScriptModule{},
	}

	p := pipeline.Pipeline{
		Name:          "Uninstall Terminus",
		Modules:       m,
		Runtime:       &runtime,
		PipelineCache: cache.NewCache(),
	}

	p.PipelineCache.Set("proxy", o.Proxy)
	p.PipelineCache.Set("kube_type", o.KubeType)

	return p.Start()
}

func Install() {

}

func Restore() {

}
