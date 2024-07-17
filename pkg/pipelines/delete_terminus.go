package pipelines

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/bootstrap/precheck"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/pipeline"
	"bytetrade.io/web3os/installer/pkg/phase/cluster"
)

func UninstallTerminusPipeline() error {
	var args = common.Argument{
		KubernetesVersion: cluster.GetCurrentKubeVersion(),
		DeleteCRI:         true,
	}

	runtime, err := common.NewKubeRuntime(common.AllInOne, args)
	if err != nil {
		return err
	}

	m := []module.Module{
		&precheck.GetStorageKeyModule{},
	}
	var kubeModules []module.Module
	switch runtime.Cluster.Kubernetes.Type {
	case common.K3s:
		kubeModules = cluster.NewK3sDeleteClusterPhase(runtime)
	case common.Kubernetes:
		kubeModules = cluster.NewK8sDeleteClusterPhase(runtime)
	default:
		return fmt.Errorf("invalid kubernetes type: %s", runtime.Cluster.Kubernetes.Type)
	}

	m = append(m, kubeModules...)

	p := pipeline.Pipeline{
		Name:    "Delete Terminus",
		Runtime: runtime,
		Modules: m,
	}

	if err := p.Start(); err != nil {
		logger.Errorf("delete terminus failed: %v", err)
		return err
	}

	return nil

}
