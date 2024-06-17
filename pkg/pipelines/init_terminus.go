package pipelines

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/pipeline"
	"bytetrade.io/web3os/installer/pkg/install"
	"bytetrade.io/web3os/installer/pkg/packages"
)

func InstallTerminusPipeline(args common.Argument) error {
	var loaderType string
	if args.FilePath != "" {
		loaderType = common.File
	} else {
		loaderType = common.AllInOne
	}

	runtime, err := common.NewKubeRuntime(loaderType, args)
	if err != nil {
		return err
	}

	m := []module.Module{
		&packages.PackagesModule{},
		&install.InstallTerminusModule{},
	}

	p := pipeline.Pipeline{
		Name:    "InstallTerminusPipeline",
		Modules: m,
		Runtime: runtime,
	}

	go p.Start()

	return nil
}
