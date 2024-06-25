package shell

import "bytetrade.io/web3os/installer/pkg/core/module"

type ExecuteShell struct {
	module.BaseTaskModule
}

func (m *ExecuteShell) GetName() string {
	return "ExecuteShell"
}

func (m *ExecuteShell) Init() {
	m.Name = "ExecuteShell"
	m.Desc = "exec shell file"

}
