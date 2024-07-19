package cluster

import (
	"bytetrade.io/web3os/installer/pkg/bootstrap/os"
	"bytetrade.io/web3os/installer/pkg/bootstrap/patch"
	"bytetrade.io/web3os/installer/pkg/bootstrap/precheck"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/pipeline"
	"bytetrade.io/web3os/installer/pkg/storage"
)

func CreateTerminus(args common.Argument, runtime *common.KubeRuntime) *pipeline.Pipeline {
	var storageVendor = args.Storage.StorageVendor
	var storageType = args.Storage.StorageType

	m := []module.Module{
		&precheck.GreetingsModule{},
		&precheck.GetSysInfoModel{},
		&precheck.PreCheckOsModule{}, // * 对应 precheck_os()
		&patch.InstallDepsModule{},   // * 对应 install_deps
		&os.ConfigSystemModule{},     // * 对应 config_system
		&storage.InitStorageModule{Skip: storageVendor != "true"},
		&storage.InstallMinioModule{Skip: storageType != "minio"},
		&storage.InstallRedisModule{},
		// &storage.InstallJuiceFsModule{},
	}

	var kubeModules []module.Module
	if runtime.Cluster.Kubernetes.Type == common.K3s {
		kubeModules = NewK3sCreateClusterPhase(runtime)
	} else {
		kubeModules = NewCreateClusterPhase(runtime)
	}

	if kubeModules == nil {
	}

	// m = append(m, kubeModules...)  // ! 暂时取消，主要测试 storage 的安装

	return &pipeline.Pipeline{
		Name:    "Install Terminus",
		Modules: m,
		Runtime: runtime,
	}
}
