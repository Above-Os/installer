package pipelines

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/pipeline"
	"bytetrade.io/web3os/installer/pkg/install"
	"bytetrade.io/web3os/installer/pkg/packages"
	"bytetrade.io/web3os/installer/pkg/storage"
)

// + 测试函数，测试下载 full 包，并执行安装
func InstallTerminusPipeline(args common.Argument) error {
	runtime, err := common.NewKubeRuntime(common.AllInOne, args)
	if err != nil {
		return err
	}

	m := []module.Module{
		&storage.SaveInstallConfigModule{},
		&packages.PackagesModule{},
		&install.InstallTerminusModule{},
	}

	p := pipeline.Pipeline{
		Name:    "Install Terminus",
		Modules: m,
		Runtime: runtime,
	}

	go p.Start()

	return nil
}
