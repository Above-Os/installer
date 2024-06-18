package mock

import (
	"bytetrade.io/web3os/installer/pkg/bootstrap/precheck"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/pipeline"
)

func NewGreetingsPipeline(runtime *common.KubeRuntime) error {
	m := []module.Module{
		&precheck.GreetingsModule{},
	}

	p := pipeline.Pipeline{
		Name:    "GreetingsPipeline",
		Modules: m,
		Runtime: runtime,
	}

	go p.Start()

	return nil
}

func Greetings(args common.Argument) error {
	runtime, err := common.NewKubeRuntime(common.AllInOne, args)
	if err != nil {
		return err
	}

	if err := NewGreetingsPipeline(runtime); err != nil {
		return err
	}
	return nil
}
