package startup

import (
	"bytetrade.io/web3os/installer/pkg/bootstrap/precheck"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/pipeline"
)

func GetMachineInfo() error {
	// runtime, err := common.NewLocalRuntime(false, false)
	runtime, err := common.NewKubeRuntime(common.AllInOne, common.Argument{})
	if err != nil {
		return err
	}

	m := []module.Module{
		&precheck.GetSysInfoModel{},
	}

	p := pipeline.Pipeline{
		Name:    "Startup",
		Modules: m,
		Runtime: runtime,
	}

	return p.Start()
}
