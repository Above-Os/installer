package binaries

import (
	"strings"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/prepare"
	"bytetrade.io/web3os/installer/pkg/core/util"
)

type Ubuntu24AppArmorCheck struct {
	prepare.BasePrepare
}

func (p *Ubuntu24AppArmorCheck) PreCheck(runtime connector.Runtime) (bool, error) {
	if constants.OsType != common.Linux || constants.OsPlatform != common.Ubuntu {
		return true, nil
	}

	if !strings.HasPrefix(constants.OsVersion, "24.") {
		return true, nil
	}

	cmd := "apparmor_parser --version"
	stdout, _, err := util.Exec(cmd, true, true)
	if err != nil {
		logger.Errorf("check apparmor version error %v", err)
		return true, nil
	}

	if strings.Index(stdout, "4.0.1") < 0 {
		return false, nil // need to install
	}

	return true, nil
}
