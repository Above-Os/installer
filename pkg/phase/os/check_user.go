package os

import (
	_ "bytetrade.io/web3os/installer/pkg/bootstrap/runuser"
	"bytetrade.io/web3os/installer/pkg/common"
)

func CheckCurrentUserPipeline() error {
	_, err := common.NewLocalRuntime(false, false)
	if err != nil {
		return err
	}

	return nil

	// m := []module.Module{
	// 	&runuser.RunUserModule{},
	// }

	// p := pipeline.Pipeline{
	// 	Name:    "CheckCurrentUserPipeline",
	// 	Modules: m,
	// }

	// return p.Start()
}
