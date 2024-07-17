package storage

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/task"
)

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
