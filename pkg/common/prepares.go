package common

import (
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/prepare"
)

type Stop struct {
	prepare.BasePrepare
}

func (p *Stop) PreCheck(runtime connector.Runtime) (bool, error) {
	return true, nil
	// return false, fmt.Errorf("STOP !!!!!!")
}