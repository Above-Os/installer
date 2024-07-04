package patch

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/prepare"
)

type CheckDepsPrepare struct {
	prepare.BasePrepare
	Command string
}

func (p *CheckDepsPrepare) PreCheck(runtime connector.Runtime) (bool, error) {
	switch constants.OsPlatform {
	case common.Ubuntu, common.Debian, common.Raspbian, common.CentOs, common.Fedora, common.RHEl:
		return true, nil
	}

	if _, err := runtime.GetRunner().Host.GetCommand(p.Command); err == nil {
		return true, nil
	}

	return false, nil
}
