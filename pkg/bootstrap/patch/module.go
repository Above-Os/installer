package patch

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/task"
)

type InstallDepsModule struct {
	module.BaseTaskModule
}

func (m *InstallDepsModule) Init() {
	m.Name = "InstallDeps"

	patchOs := &task.LocalTask{
		Name:   "PatchOs",
		Action: new(PatchTask),
		Retry:  0,
	}

	installSocat := &task.LocalTask{
		Name:    "InstallSocat",
		Prepare: &CheckDepsPrepare{Command: common.CommandSocat},
		Action:  new(SocatTask),
	}

	installConntrack := &task.LocalTask{
		Name:    "InstallConntrack",
		Prepare: &CheckDepsPrepare{Command: common.CommandConntrack},
		Action:  new(ConntrackTask),
	}

	m.Tasks = []task.Interface{
		patchOs,
		installSocat,
		installConntrack,
	}
}
