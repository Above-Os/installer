package pipelines

import (
	"bytetrade.io/web3os/installer/pkg/bootstrap/precheck"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/pipeline"
)

// todo 安装 Terminus
// todo 这里先考虑从 mini 包进行安装
func NewCreateInstallerPipeline(runtime *common.KubeRuntime) error {
	// precheck/GetSysInfoModel    程序启动时就已经执行了，这里不再执行
	// binaries/patch ubuntu24 apparmor

	m := []module.Module{
		&precheck.TerminusGreetingsModule{},
		&precheck.PreCheckOsModule{},
	}

	p := pipeline.Pipeline{
		Name:    "CreateInstallPipeline",
		Modules: m,
		Runtime: runtime,
	}

	go p.Start()

	return nil
}
