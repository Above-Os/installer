package pipelines

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/pipeline"
	"bytetrade.io/web3os/installer/pkg/install"
	"bytetrade.io/web3os/installer/pkg/packages"
	"bytetrade.io/web3os/installer/pkg/scripts"
	"bytetrade.io/web3os/installer/pkg/storage"
)

func InstallTerminusPipeline(args common.Argument) error {
	runtime, err := common.NewKubeRuntime(common.AllInOne, args) // 后续拆解 install_cmd.sh，会用到 KubeRuntime
	if err != nil {
		return err
	}

	m := []module.Module{
		&storage.SaveInstallConfigModule{},
		&packages.PackagesModule{},
		&scripts.CopyUninstallScriptModule{},
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
