package system

import (
	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/core/action"
	corecommon "bytetrade.io/web3os/installer/pkg/core/common"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/task"
)

// ~ InstallDepsModule
type InstallDepsModule struct {
	module.BaseTaskModule
}

func (m *InstallDepsModule) GetName() string {
	return "InstallDepsModule"
}

func (m *InstallDepsModule) Init() {
	m.Name = "InstallDepsModule"
	m.Desc = "Install dependencies"

	installDeps := &task.LocalTask{
		Name: "InstallDeps",
		Desc: "Install dependencies",
		Action: &action.Script{
			Name: "installDeps",
			File: corecommon.GreetingShell,
			Args: []string{constants.OsPlatform},
		},
	}

	m.Tasks = []task.Interface{
		installDeps,
	}
}
