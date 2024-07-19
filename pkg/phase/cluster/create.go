package cluster

import (
	"bytetrade.io/web3os/installer/pkg/bootstrap/os"
	"bytetrade.io/web3os/installer/pkg/bootstrap/patch"
	"bytetrade.io/web3os/installer/pkg/bootstrap/precheck"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/pipeline"
)

func CreateTerminus(args common.Argument, runtime *common.KubeRuntime) *pipeline.Pipeline {
	m := []module.Module{
		&precheck.GreetingsModule{},
		&precheck.GetSysInfoModel{},
		&precheck.PreCheckOsModule{}, // * 对应 precheck_os()
		&patch.InstallDepsModule{},   // * 对应 install_deps
		&os.ConfigSystemModule{},     // * 对应 config_system
		// todo storage
	}

	var kubeModules []module.Module
	if runtime.Cluster.Kubernetes.Type == common.K3s {
		kubeModules = NewK3sCreateClusterPhase(runtime)
	} else {
		kubeModules = NewCreateClusterPhase(runtime)
	}

	m = append(m, kubeModules...)

	return &pipeline.Pipeline{
		Name:    "Install Terminus",
		Modules: m,
		Runtime: runtime,
	}
}
