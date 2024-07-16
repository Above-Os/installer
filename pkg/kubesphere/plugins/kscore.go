package plugins

import (
	"context"
	"fmt"
	"path"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/prepare"
	"bytetrade.io/web3os/installer/pkg/core/task"
	"bytetrade.io/web3os/installer/pkg/utils"
	ctrl "sigs.k8s.io/controller-runtime"
)

type CreateKsCore struct {
	common.KubeAction
}

func (t *CreateKsCore) Execute(runtime connector.Runtime) error {
	masterNumIf, ok := t.PipelineCache.Get(common.CacheMasterNum)
	if !ok || masterNumIf == nil {
		return fmt.Errorf("failed to get master num")
	}
	masterNum := masterNumIf.(int64)

	config, err := ctrl.GetConfig()
	if err != nil {
		return err
	}

	var appKsCoreName = common.ChartNameKsCore
	var appPath = path.Join(runtime.GetFilesDir(), "apps", appKsCoreName)

	actionConfig, settings, err := utils.InitConfig(config, common.NamespaceKubesphereSystem)
	if err != nil {
		return err
	}

	var values = make(map[string]interface{})
	values["Release"] = map[string]string{
		"Namespace":    common.NamespaceKubesphereSystem,
		"ReplicaCount": fmt.Sprintf("%s", masterNum),
	}
	if err := utils.InstallCharts(context.Background(), actionConfig, settings, appKsCoreName,
		appPath, "", common.NamespaceKubesphereSystem, values); err != nil {
		logger.Errorf("failed to install %s chart: %v", appKsCoreName, err)
		return err
	}

	return nil
}

// ~ DeployKsCoreModule
type DeployKsCoreModule struct {
	common.KubeModule
}

func (m *DeployKsCoreModule) Init() {
	m.Name = "DeployKsCore"

	createKsCore := &task.RemoteTask{
		Name:  "CreateKsCore",
		Hosts: m.Runtime.GetHostsByRole(common.Master),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
			new(NotEqualDesiredVersion),
		},
		Action:   new(CreateKsCore),
		Parallel: false,
		Retry:    0,
	}

	m.Tasks = []task.Interface{
		createKsCore,
	}
}
