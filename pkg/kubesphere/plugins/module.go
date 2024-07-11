package plugins

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/prepare"
	"bytetrade.io/web3os/installer/pkg/core/task"
	"bytetrade.io/web3os/installer/pkg/kubernetes"
)

type DeployKsPluginsModule struct {
	common.KubeModule
}

func (d *DeployKsPluginsModule) Init() {
	d.Name = "DeployKsPluginsModule"

	newNamespace := &task.RemoteTask{
		Name:  "CreateKsPluginsNamespace",
		Hosts: d.Runtime.GetHostsByRole(common.Master),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
			new(NotEqualDesiredVersion),
		},
		Action:   new(InitNamespace),
		Parallel: false,
	}

	createSnapshotController := &task.RemoteTask{
		Name:  "CreateSnapshotController",
		Hosts: d.Runtime.GetHostsByRole(common.Master),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
			new(NotEqualDesiredVersion),
			&kubernetes.GetKubeletVersion{
				CommandDelete: false,
			},
		},
		Action:   new(DeploySnapshotController),
		Parallel: false,
	}

	d.Tasks = []task.Interface{
		newNamespace,
		createSnapshotController,
	}
}
