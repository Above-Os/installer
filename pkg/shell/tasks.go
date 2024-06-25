package shell

import (
	"bytetrade.io/web3os/installer/pkg/core/action"
	"bytetrade.io/web3os/installer/pkg/core/connector"
)

type ExecShellTask struct {
	action.BaseAction
}

func (t *ExecShellTask) Execute(runtime connector.Runtime) error {
	return nil
}
