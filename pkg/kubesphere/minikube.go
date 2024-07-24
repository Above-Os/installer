package kubesphere

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/task"
	"bytetrade.io/web3os/installer/pkg/core/util"
)

// ~ CopyFiles
type CopyFiles struct {
	common.KubeAction
}

func (t *CopyFiles) Execute(runtime connector.Runtime) error {
	var cmd = "/usr/local/bin/kubectl get pods -A"
	if _, err := runtime.GetRunner().Host.Cmd(cmd, false, true); err != nil {
		return err
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

	checkMacCommandExists := &task.RemoteTask{
		Name:     "CheckMiniKubeExists",
		Hosts:    m.Runtime.GetHostsByRole(common.Master),
		Prepare:  new(common.OnlyFirstMaster),
		Action:   new(CheckMacCommandExists),
		Parallel: false,
	}

	copyFiles := &task.RemoteTask{
		Name:     "CopyFiles",
		Hosts:    m.Runtime.GetHostsByRole(common.Master),
		Prepare:  new(common.OnlyFirstMaster),
		Action:   new(CopyFiles),
		Parallel: false,
	}

	m.Tasks = []task.Interface{
		checkMacCommandExists,
		copyFiles,
	}
}
