package runuser

import "bytetrade.io/web3os/installer/pkg/core/module"

type RunUserModule struct {
	module.BaseModule
}

func (m *RunUserModule) Init() {
	m.Name = "RunUserModule"
	m.Desc = "running user check"

	m.PostHook = []module.PostHookInterface{
		&RunUserCheckHook{},
	}
}
