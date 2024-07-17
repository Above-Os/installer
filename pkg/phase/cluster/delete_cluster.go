package cluster

import (
	"fmt"
	"strings"

	"bytetrade.io/web3os/installer/pkg/bootstrap/os"
	"bytetrade.io/web3os/installer/pkg/bootstrap/precheck"
	"bytetrade.io/web3os/installer/pkg/certs"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/container"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"bytetrade.io/web3os/installer/pkg/k3s"
	"bytetrade.io/web3os/installer/pkg/kubernetes"
	"bytetrade.io/web3os/installer/pkg/loadbalancer"
	"bytetrade.io/web3os/installer/pkg/utils"
)

func GetCurrentKubeVersion() string {
	var fileVersion string
	var f = "/etc/kke/version"
	if ok := utils.IsExist(f); ok {
		stdout, _, err := util.Exec("awk -F '=' '/KUBE/{printf \"%s\",$2}' "+f, false, true)
		if err != nil {
			logger.Errorf("get kube version error %v", err)
		}
		if stdout != "" {
			fileVersion = stdout
			constants.InstalledKubeVersion = fileVersion
		}
	}

	kubectl, _, err := util.Exec(fmt.Sprintf("command -v %s", common.CommandKubectl), false, false)
	if err != nil {
		goto SKIP
	}

	if kubectl != "" {
		stdout, _, err := util.Exec(fmt.Sprintf("%s get nodes -o jsonpath='{.items[0].status.nodeInfo.kubeletVersion}'", kubectl), false, false)
		if err != nil {
			goto SKIP
		}
		if stdout != "" {
			if strings.Contains(stdout, "+k3s1") {
				stdout = strings.ReplaceAll(stdout, "+k3s1", "-k3s")
			} else if strings.Contains(stdout, "+k3s2") {
				stdout = strings.ReplaceAll(stdout, "+k3s2", "-k3s")
			}
		}
		if fileVersion != "" {
			if stdout != "" && strings.Contains(stdout, fileVersion) {
				constants.InstalledKubeVersion = fileVersion
			}
		} else {
			if stdout != "" {
				constants.InstalledKubeVersion = stdout
			}
		}
	}
	goto SKIP

SKIP:
	if constants.InstalledKubeVersion != "" {
		fmt.Printf("KUBE: version: %s\n", constants.InstalledKubeVersion)
	}

	return constants.InstalledKubeVersion
}

func NewK8sDeleteClusterPhase(runtime *common.KubeRuntime) []module.Module {
	return []module.Module{
		&precheck.GreetingsModule{},
		&kubernetes.ResetClusterModule{},
		&container.UninstallContainerModule{Skip: !runtime.Arg.DeleteCRI},
		&os.ClearOSEnvironmentModule{},
		&certs.UninstallAutoRenewCertsModule{},
		&loadbalancer.DeleteVIPModule{Skip: !runtime.Cluster.ControlPlaneEndpoint.IsInternalLBEnabledVip()},
	}
}

func NewK3sDeleteClusterPhase(runtime *common.KubeRuntime) []module.Module {
	return []module.Module{
		&precheck.GreetingsModule{},
		&k3s.DeleteClusterModule{},
		&os.ClearOSEnvironmentModule{},
		&certs.UninstallAutoRenewCertsModule{},
		&loadbalancer.DeleteVIPModule{Skip: !runtime.Cluster.ControlPlaneEndpoint.IsInternalLBEnabledVip()},
	}
}
