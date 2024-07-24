package pipelines

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"

	ctrl "bytetrade.io/web3os/installer/controllers"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/phase/cluster"
)

func CliInitializeTerminusPipeline(kubeType string) error {
	arg := common.Argument{
		KsEnable:         true,
		KsVersion:        common.DefaultKubeSphereVersion,
		InstallPackages:  false,
		SKipPushImages:   false,
		ContainerManager: common.Containerd,
	}
	switch kubeType {
	case common.K3s:
		arg.KubernetesVersion = common.DefaultK3sVersion
	case common.K8s:
		arg.KubernetesVersion = common.DefaultK8sVersion
	}

	runtime, err := common.NewKubeRuntime(common.AllInOne, arg)
	if err != nil {
		return nil
	}
	var p = cluster.InitKube(arg, runtime)
	if err := p.Start(); err != nil {
		return fmt.Errorf("create terminus error %v", err)
	}

	if runtime.Cluster.KubeSphere.Enabled {

		fmt.Print(`Installation is complete.

Please check the result using the command:

	kubectl logs -n kubesphere-system $(kubectl get pod -n kubesphere-system -l 'app in (ks-install, ks-installer)' -o jsonpath='{.items[0].metadata.name}') -f   

`)
	} else {
		fmt.Print(`Installation is complete.

Please check the result using the command:
		
	kubectl get pod -A

`)

	}

	if runtime.Arg.InCluster {
		if err := ctrl.UpdateStatus(runtime); err != nil {
			logger.Errorf("failed to update status: %v", err)
			return err
		}
		kkConfigPath := filepath.Join(runtime.GetWorkDir(), fmt.Sprintf("config-%s", runtime.ObjName))
		if config, err := ioutil.ReadFile(kkConfigPath); err != nil {
			logger.Errorf("failed to read kubeconfig: %v", err)
			return err
		} else {
			runtime.Kubeconfig = base64.StdEncoding.EncodeToString(config)
			if err := ctrl.UpdateKubeSphereCluster(runtime); err != nil {
				logger.Errorf("failed to update kubesphere cluster: %v", err)
				return err
			}
			if err := ctrl.SaveKubeConfig(runtime); err != nil {
				logger.Errorf("failed to save kubeconfig: %v", err)
				return err
			}
		}
	}

	return nil
}
