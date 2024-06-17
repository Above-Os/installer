package install

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"github.com/pkg/errors"
)

type Terminus struct {
	common.KubeAction
}

func (a *Terminus) Execute(runtime connector.Runtime) error {
	installCmd := "bash /home/zhaoyu/install-wizard/install_cmd.sh"
	if _, err := runtime.GetRunner().SudoCmd(installCmd, false); err != nil {
		return errors.Wrapf(errors.WithStack(err), "install terminus failed")
	}
	return nil
}
