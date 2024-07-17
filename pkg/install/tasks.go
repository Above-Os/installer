package install

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
)

// ~ WizardTask
type WizardTask struct {
	common.KubeAction
}

// todo 这里增加一个检查 terminus 是否安装完成的检查，需要重启 wizard 这容器
func (t *WizardTask) Execute(runtime connector.Runtime) error {
	return nil
}
