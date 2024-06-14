package common

import "bytetrade.io/web3os/installer/pkg/core/connector"

type PackageRuntime struct {
	connector.InstallerPackageRuntime
	Arg PackageArgument
}

type PackageArgument struct {
}

func NewPackageRuntime(arg PackageArgument) (*PackageRuntime, error) {
	installer := connector.NewInstallerPackageRuntime("", connector.NewDialer(), true, false)
	r := &PackageRuntime{
		Arg: arg,
	}
	r.InstallerPackageRuntime = installer
	return r, nil
}

// Copy is used to create a copy for Runtime.
func (k *PackageRuntime) Copy() connector.PackageRuntime {
	runtime := *k
	return &runtime
}
