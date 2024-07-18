package storage

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/prepare"
	"bytetrade.io/web3os/installer/pkg/core/task"
)

// ~ RemoveStorage
type RemoveStorage struct {
	common.KubeModule
}

func (m *RemoveStorage) Init() {
	m.Name = "RemoveStorage"

	stopJuiceFS := &task.RemoteTask{
		Name:  "StopJuiceFS",
		Hosts: m.Runtime.GetHostsByRole(common.Master),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
		},
		Action:   new(StopJuiceFS),
		Parallel: false,
		Retry:    0,
	}

	stopMinio := &task.RemoteTask{
		Name:  "StopMinio",
		Hosts: m.Runtime.GetHostsByRole(common.Master),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
		},
		Action:   new(StopMinio),
		Parallel: false,
		Retry:    0,
	}

	stopMinioOperator := &task.RemoteTask{
		Name:  "StopMinioOperator",
		Hosts: m.Runtime.GetHostsByRole(common.Master),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
		},
		Action:   new(StopMinioOperator),
		Parallel: false,
		Retry:    0,
	}

	stopRedis := &task.RemoteTask{
		Name:  "StopRedis",
		Hosts: m.Runtime.GetHostsByRole(common.Master),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
		},
		Action:   new(StopRedis),
		Parallel: false,
		Retry:    0,
	}

	removeTerminusFiles := &task.RemoteTask{
		Name:  "RemoveTerminusFiles",
		Hosts: m.Runtime.GetHostsByRole(common.Master),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
		},
		Action:   new(RemoveTerminusFiles),
		Parallel: false,
		Retry:    0,
	}

	m.Tasks = []task.Interface{
		stopJuiceFS,
		stopMinio,
		stopMinioOperator,
		stopRedis,
		removeTerminusFiles,
	}
}

// ~ SaveInstallConfigModule
type SaveInstallConfigModule struct {
	common.KubeModule
}

func (m *SaveInstallConfigModule) Init() {
	m.Name = "SaveInstallConfig"

	save := &task.RemoteTask{
		Name:     "Save",
		Hosts:    m.Runtime.GetAllHosts(),
		Action:   &SaveInstallConfigTask{},
		Parallel: false,
		Retry:    0,
	}

	m.Tasks = []task.Interface{save}
}
