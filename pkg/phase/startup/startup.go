package startup

import (
	"bytetrade.io/web3os/installer/pkg/bootstrap/precheck"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/pipeline"
)

func GetMachineInfo() error {
	runtime, err := common.NewLocalRuntime(false, false)
	if err != nil {
		return err
	}

	m := []module.Module{
		&precheck.GetSysInfoModel{},
	}

	p := pipeline.Pipeline{
		Name:    "CreateStartupPipeline",
		Modules: m,
		Runtime: &runtime,
	}

	return p.Start()
}
