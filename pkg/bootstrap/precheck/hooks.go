package precheck

import (
	"bytetrade.io/web3os/installer/pkg/core/ending"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/module"

	gphost "github.com/shirou/gopsutil/v4/host"
)

type GatherHook struct {
	Module module.Module
	Result *ending.ModuleResult
}

func (h *GatherHook) Init(module module.Module, result *ending.ModuleResult) {

}

func (h *GatherHook) Try() error {
	gpis, err := gphost.Info()
	if err != nil {
		return err
	}
	logger.Debugf("host info os: %s, hostname: %s, platform: %s, arch: %s, version: %s", gpis.OS, gpis.Hostname, gpis.Platform, gpis.KernelArch, gpis.KernelVersion)
	return nil
}

func (h *GatherHook) Catch(err error) error {
	return nil
}

func (h *GatherHook) Finally() {

}
