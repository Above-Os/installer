package dm

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/task"
)

type DmModule struct {
	common.KubeModule
}

func (m *DmModule) Init() {
	m.Name = "dm"
	m.Desc = "dm module"

	t1 := &task.RemoteTask{
		Name:    "Task1",
		Hosts:   m.Runtime.GetAllHosts(),
		Prepare: &P1{},
		Action:  &Task1{},
	}

	t2 := &task.RemoteTask{
		Name:   "Task2",
		Hosts:  m.Runtime.GetAllHosts(),
		Action: &Task2{},
	}

	t3 := &task.RemoteTask{
		Name:   "Task3",
		Hosts:  m.Runtime.GetAllHosts(),
		Action: &Task3{},
	}

	m.Tasks = []task.Interface{
		t1,
		t2,
		t3,
	}
}
