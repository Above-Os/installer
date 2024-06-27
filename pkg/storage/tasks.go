package storage

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
)

// ~ task SaveInstallConfigTask
type SaveInstallConfigTask struct {
	common.KubeAction
}

func (t *SaveInstallConfigTask) Execute(runtime connector.Runtime) error {
	logger.Debug("[action] SaveInstallConfigTask")
	fmt.Printf("---task--- %+v\n", t.KubeConf.Arg.Request)
	return t.KubeConf.Arg.Provider.Ping()
}
