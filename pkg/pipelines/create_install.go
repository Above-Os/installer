package pipelines

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"

	ctrl "bytetrade.io/web3os/installer/controllers"
	"bytetrade.io/web3os/installer/pkg/bootstrap/os"
	"bytetrade.io/web3os/installer/pkg/bootstrap/patch"
	"bytetrade.io/web3os/installer/pkg/bootstrap/precheck"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/pipeline"
	"bytetrade.io/web3os/installer/pkg/phase/cluster"
)

// todo 安装 Terminus
// todo 这里先考虑从 mini 包进行安装
// + 测试
func NewCreateInstallerPipeline(runtime *common.KubeRuntime) error {
	// precheck/GetSysInfoModel    程序启动时就已经执行了，这里不再执行
	// binaries/patch ubuntu24 apparmor

	m := []module.Module{
		&precheck.TerminusGreetingsModule{},
		&precheck.PreCheckOsModule{}, // * 对应 precheck_os()
		&patch.InstallDepsModule{},   // * 对应 install_deps
		&os.ConfigSystemModule{},     // * 对应 config_system
	}

	modules := cluster.NewK3sCreateClusterPhase(runtime)
	m = append(m, modules...)

	p := pipeline.Pipeline{
		Name:    "CreateInstallPipeline",
		Modules: m,
		Runtime: runtime,
	}

	go func() {
		if err := p.Start(); err != nil {
			logger.Errorf("install k3s failed %v", err)
			return
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
				return
			}
			kkConfigPath := filepath.Join(runtime.GetWorkDir(), fmt.Sprintf("config-%s", runtime.ObjName))
			if config, err := ioutil.ReadFile(kkConfigPath); err != nil {
				logger.Errorf("failed to read kubeconfig: %v", err)
				return
			} else {
				runtime.Kubeconfig = base64.StdEncoding.EncodeToString(config)
				if err := ctrl.UpdateKubeSphereCluster(runtime); err != nil {
					logger.Errorf("failed to update kubesphere cluster: %v", err)
					return
				}
				if err := ctrl.SaveKubeConfig(runtime); err != nil {
					logger.Errorf("failed to save kubeconfig: %v", err)
					return
				}
			}
		}
	}()

	return nil
}
