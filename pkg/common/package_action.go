package common

import (
	"bytetrade.io/web3os/installer/pkg/core/action"
	"bytetrade.io/web3os/installer/pkg/core/connector"
)

type PackageAction struct {
	action.InstallerBaseAction
}

func (k *PackageAction) AutoAssert(runtime connector.InstallerPackageRuntime) {

}
