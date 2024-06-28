package binaries

import (
	"strings"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/prepare"
	"bytetrade.io/web3os/installer/pkg/core/util"
)

type Ubuntu24AppArmorCheck struct {
	prepare.BasePrepare
}

func (p *Ubuntu24AppArmorCheck) PreCheck(runtime connector.Runtime) (bool, error) {
	if constants.OsType != common.Linux {
		return true, nil
	}

	if !strings.HasPrefix(constants.OsVersion, "24.") {
		return true, nil
	}

	cmd := "apparmor_parser --version"
	stdout, _, err := util.Exec(cmd, true, false)
	if err != nil {
		return true, err
	}

	if strings.Index(stdout, "4.0.1") < 0 {
		// 这里是表示需要安装 apparmor
		return false, nil
	}

	return true, nil
}
