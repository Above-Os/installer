package kk

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"github.com/pkg/errors"
)

// ~ test for install kk
type ChmodKk struct {
	common.KubeAction
}

func (a *ChmodKk) Execute(runtime connector.Runtime) error {
	fmt.Println("[action] ChmodKk")
	if _, err := runtime.GetRunner().SudoCmd("chmod +x /tmp/install_log/kk", false); err != nil {
		return errors.Wrapf(errors.WithStack(err), "chmod kk failed")
	}
	return nil
}

type ExecuteKk struct {
	common.KubeAction
}

func (a *ExecuteKk) Execute(runtime connector.Runtime) error {
	fmt.Println("[action] ExecuteKk")
	// kk 的安装走的是脚本
	installCmd := "/tmp/install_log/kk create cluster --with-kubernetes v1.21.4-k3s --with-kubesphere v3.3.0 --container-manager containerd "
	if _, err := runtime.GetRunner().SudoCmd(installCmd, false); err != nil {
		return errors.Wrapf(errors.WithStack(err), "install kk failed")
	}

	return nil

}
