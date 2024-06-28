package storage

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/task"
)

type SaveInstallConfigModule struct {
	common.KubeModule
}

func (m *SaveInstallConfigModule) Init() {
	m.Name = "SaveInstallConfigModule"
	m.Desc = "SaveInstallConfigModule"

	save := &task.LocalTask{
		Name:   "SaveInstallConfig",
		Desc:   "SaveInstallConfig",
		Action: &SaveInstallConfigTask{},
		Retry:  0,
	}

	m.Tasks = []task.Interface{save}
}
