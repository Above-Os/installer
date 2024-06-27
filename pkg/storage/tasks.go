package storage

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
)

// ~ task SaveInstallConfigTask
type SaveInstallConfigTask struct {
	common.KubeAction
}

func (t *SaveInstallConfigTask) GetName() string {
	return "SaveInstallConfigTask"
}

func (t *SaveInstallConfigTask) Execute(runtime connector.Runtime) error {
	fmt.Printf("---task--- %+v\n", t.KubeConf.Arg.Request)
	return t.KubeConf.Arg.Provider.Ping()
}
