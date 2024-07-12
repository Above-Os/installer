package plugins

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/prepare"
	"bytetrade.io/web3os/installer/pkg/core/task"
)

type CopyEmbed struct {
	common.KubeModule
}

func (t *CopyEmbed) Init() {
	t.Name = "CopyEmbed"

	copyEmbed := &task.LocalTask{
		Name:   "CopyEmbedFiles",
		Action: new(CopyEmbedFiles),
	}

	t.Tasks = []task.Interface{
		copyEmbed,
	}
}

type DeployKsPluginsModule struct {
	common.KubeModule
}

func (t *DeployKsPluginsModule) Init() {
	t.Name = "DeployKsPlugins"

	initNs := &task.RemoteTask{
		Name:  "InitKsNamespace",
		Hosts: t.Runtime.GetHostsByRole(common.Master),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
			new(NotEqualDesiredVersion),
		},
		Action:   new(InitNamespace),
		Parallel: true,
	}

	// checkMasterNum := &task.RemoteTask{
	// 	Name:  "CheckMasterNum",
	// 	Hosts: t.Runtime.GetHostsByRole(common.Master),
	// 	Prepare: &prepare.PrepareCollection{
	// 		new(common.OnlyFirstMaster),
	// 		new(NotEqualDesiredVersion),
	// 	},
	// 	Action:   new(CheckMasterNum),
	// 	Parallel: true,
	// }

	t.Tasks = []task.Interface{
		initNs,
		// checkMasterNum,
	}
}

// +

type DebugModule struct {
	common.KubeModule
}

func (m *DebugModule) Init() {
	m.Name = "Debug"

	patchRedis := &task.RemoteTask{
		Name:  "PatchRedis",
		Hosts: m.Runtime.GetHostsByRole(common.ETCD),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
			new(NotEqualDesiredVersion),
		},
		Action:   new(PatchRedisStatus),
		Parallel: true,
	}

	m.Tasks = []task.Interface{
		patchRedis,
	}
}
