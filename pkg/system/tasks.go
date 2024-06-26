package system

import (
	"bytetrade.io/web3os/installer/pkg/core/action"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
)

type InstallDeps struct {
	action.BaseAction
}

func (i *InstallDeps) Execute(runtime connector.Runtime) error {
	logger.Debug("[action] InstallDeps")
	return nil
}
