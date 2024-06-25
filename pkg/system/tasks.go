package system

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/util"
)

type InstallDeps struct {
	common.KubeAction
}

func (t *InstallDeps) Execute(runtime connector.Runtime) error {
	logger.Debug("[action] DepsUpdate")
	var pre_reqs = "apt-transport-https ca-certificates curl"
	var cmd string
	var pkgManager string = "yum"
	switch constants.OsPlatform {
	case common.Ubuntu, common.Debian, common.Raspbian:
		_, _, err := util.Exec("command -v gpg > /dev/null", false)
		if err != nil {
			pre_reqs = fmt.Sprintf("%s gnupg", pre_reqs)
		}
		_, _, err = util.Exec("apt-get update -qq >/dev/null", false)
		if err != nil {
			return err
		}
		_, _, err = util.Exec(fmt.Sprintf("DEBIAN_FRONTEND=noninteractive apt-get install -y -qq %s >/dev/null", pre_reqs), false)
		if err != nil {
			return err
		}
		// _,_, err =

	case common.CentOs, common.RHEl, common.Fedora:
		if constants.OsPlatform == common.Fedora {
			pkgManager = "dnf"
		}
		pkgManager = "yum"
		_, _, err := util.Exec(fmt.Sprintf("%s install -y conntrack socat httpd-tools ntpdate net-tools make gcc openssh-server", pkgManager), false)
		if err != nil {
			fmt.Println("---install deps error---", err)
			return err
		}
	default:
		// todo 单独执行其他的安装方式
		return fmt.Errorf("platform %s not support", constants.OsPlatform)
	}

	return nil
}
