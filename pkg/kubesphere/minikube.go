package kubesphere

import (
	"fmt"
	"strings"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/task"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"github.com/pkg/errors"
)

// ~ InitMinikubeNs
type InitMinikubeNs struct {
	common.KubeAction
}

func (t *InitMinikubeNs) Execute(runtime connector.Runtime) error {
	var allNs = []string{
		common.NamespaceDefault,
		common.NamespaceKubeNodeLease,
		common.NamespaceKubePublic,
		common.NamespaceKubeSystem,
		common.NamespaceKubekeySystem,
		common.NamespaceKubesphereControlsSystem,
		common.NamespaceKubesphereMonitoringFederated,
		common.NamespaceKubesphereMonitoringSystem,
		common.NamespaceKubesphereSystem,
	}

	for _, ns := range allNs {
		if stdout, err := runtime.GetRunner().SudoCmdExt(fmt.Sprintf("/usr/local/bin/kubectl create ns %s", ns), false, true); err != nil {
			if !strings.Contains(stdout, "already exists") {
				logger.Errorf("create ns %s failed: %v", ns, err)
				return errors.Wrap(errors.WithStack(err), fmt.Sprintf("create namespace %s failed: %v", ns, err))
			}
		}
	}

	return nil
}

// ~ CheckMacCommandExists
type CheckMacCommandExists struct {
	common.KubeAction
}

func (t *CheckMacCommandExists) Execute(runtime connector.Runtime) error {
	if m, err := util.GetCommand(common.CommandMinikube); err != nil || m == "" {
		return fmt.Errorf("minikube not found")
	}

	if d, err := util.GetCommand(common.CommandDocker); err != nil || d == "" {
		return fmt.Errorf("docker not found")
	}

	if h, err := util.GetCommand(common.CommandHelm); err != nil || h == "" {
		return fmt.Errorf("helm not found")
	}

	if h, err := util.GetCommand(common.CommandKubectl); err != nil || h == "" {
		return fmt.Errorf("kubectl not found")
	}
	return nil
}

// ~ DeployMiniKubeModule
type DeployMiniKubeModule struct {
	common.KubeModule
}

func (m *DeployMiniKubeModule) Init() {
	m.Name = "DeployMacOS"

	checkMacCommandExists := &task.LocalTask{
		Name:    "CheckMiniKubeExists",
		Prepare: new(common.OnlyFirstMaster),
		Action:  new(CheckMacCommandExists),
	}

	// copyFiles := &task.RemoteTask{
	// 	Name:   "CopyFiles",
	// 	Hosts:  m.Runtime.GetHostsByRole(common.Master),
	// 	Action: new(CopyFiles),
	// }

	apply := &task.LocalTask{
		Name:    "ApplyKsInstaller",
		Prepare: new(common.OnlyFirstMaster),
		Action:  new(Apply),
	}

	m.Tasks = []task.Interface{
		checkMacCommandExists,
		apply,
	}
}
