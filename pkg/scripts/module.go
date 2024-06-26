package scripts

import (
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/task"
)

// ~ CopyScriptsModule
type CopyScriptsModule struct {
	module.BaseTaskModule
}

func (m *CopyScriptsModule) GetName() string {
	return "CopyScriptsModule"
}

func (m *CopyScriptsModule) Init() {
	m.Name = "CopyScriptsModule"
	m.Desc = "Copy scripts"

	copyScripts := &task.LocalTask{
		Name:   "CopyScripts",
		Desc:   "Copy scripts",
		Action: &Copy{},
	}

	greeting := &task.LocalTask{
		Name:        "Greeting",
		Desc:        "Greeting",
		Action:      &Greeting{},
		Retry:       0,
		IgnoreError: true,
	}

	m.Tasks = []task.Interface{
		copyScripts,
		greeting,
	}
}
