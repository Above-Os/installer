package kubesphere

import (
	"fmt"
	"path/filepath"
	"strings"

	"bytetrade.io/web3os/installer/pkg/common"
	cc "bytetrade.io/web3os/installer/pkg/core/common"
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
		common.NamespaceKubekeySystem,
		common.NamespaceKubesphereSystem,
		common.NamespaceKubesphereMonitoringSystem,
	}

	for _, ns := range allNs {
		if stdout, err := runtime.GetRunner().Host.CmdExt(fmt.Sprintf("/usr/local/bin/kubectl create ns %s", ns), false, true); err != nil {
			if !strings.Contains(stdout, "already exists") {
				logger.Errorf("create ns %s failed: %v", ns, err)
				return errors.Wrap(errors.WithStack(err), fmt.Sprintf("create namespace %s failed: %v", ns, err))
			}
		}
	}

	filePath := filepath.Join(common.KubeAddonsDir, "cluster.yaml")
	if err := util.WriteFile(filePath, []byte(t.KubeConf.Cluster.KubeSphere.Configurations), cc.FileMode0755); err != nil {
		return errors.Wrap(errors.WithStack(err), fmt.Sprintf("write ks installer %s failed", filePath))
	}

	deployKubesphereCmd := fmt.Sprintf("/usr/local/bin/kubectl apply -f %s --force", filePath)
	if _, err := runtime.GetRunner().Host.CmdExt(deployKubesphereCmd, false, true); err != nil {
		return errors.Wrapf(errors.WithStack(err), "deploy %s failed", filePath)
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
		Name:   "CheckMiniKubeExists",
		Action: new(CheckMacCommandExists),
	}

	initMinikubeNs := &task.LocalTask{
		Name:   "InitMinikubeNs",
		Action: new(InitMinikubeNs),
	}

	m.Tasks = []task.Interface{
		checkMacCommandExists,
		initMinikubeNs,
	}
}
