package install

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/util"
)

type Terminus struct {
	common.KubeAction
}

func (a *Terminus) Execute(runtime connector.Runtime) error {
	fmt.Println("[action] Terminus")
	installCmd := "export TERMINUS_OS_DOMAINNAME=myterminus.com;export TERMINUS_OS_USERNAME=zhaoyu;export TERMINUS_OS_EMAIL=zhaoyu@bytetrade.io;bash /home/zhaoyu/install-wizard/install_cmd.sh"
	_, _, err := util.Exec(installCmd, false)
	if err != nil {
		return err
	}
	// runtime.GetRunner().SudoCmd(installCmd, false)
	// if _, err := runtime.GetRunner().SudoCmd(installCmd, false); err != nil {
	// 	return errors.Wrapf(errors.WithStack(err), "install terminus failed")
	// }

	return nil
}
