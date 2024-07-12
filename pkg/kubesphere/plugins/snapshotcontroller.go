package plugins

import (
	"context"
	"path"
	"time"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/prepare"
	"bytetrade.io/web3os/installer/pkg/core/task"
	"bytetrade.io/web3os/installer/pkg/kubernetes"
	"bytetrade.io/web3os/installer/pkg/utils"

	ctrl "sigs.k8s.io/controller-runtime"
)

// ~ DeploySnapshotController
type DeploySnapshotController struct {
	common.KubeAction
}

func (t *DeploySnapshotController) Execute(runtime connector.Runtime) error {
	config, err := ctrl.GetConfig()
	if err != nil {
		return err
	}

	var appName = common.ChartNameSnapshotController
	var appPath = path.Join(runtime.GetFilesDir(), "apps", appName)

	actionConfig, settings, err := utils.InitConfig(config, common.NamespaceKubeSystem)
	if err != nil {
		return err
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	var values = make(map[string]interface{})
	values["Release"] = map[string]string{
		"Namespace": common.NamespaceKubeSystem,
	}

	if err := utils.InstallCharts(ctx, actionConfig, settings, appName, appPath, "", appName, values); err != nil {
		return err
	}

	return nil
}

// ~ DeploySnapshotControllerModule
type DeploySnapshotControllerModule struct {
	common.KubeModule
}

func (d *DeploySnapshotControllerModule) Init() {
	d.Name = "DeploySnapshotController"

	createSnapshotController := &task.RemoteTask{
		Name:  "CreateSnapshotController",
		Hosts: d.Runtime.GetHostsByRole(common.Master),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
			new(NotEqualDesiredVersion),
			&kubernetes.GetKubeletVersion{
				CommandDelete: false,
			},
		},
		Action:   new(DeploySnapshotController),
		Parallel: false,
	}

	d.Tasks = []task.Interface{
		createSnapshotController,
	}
}
