package kubesphere

import (
	"encoding/base64"
	"fmt"
	"path"

	"bytetrade.io/web3os/installer/pkg/common"
	cc "bytetrade.io/web3os/installer/pkg/core/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/task"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"bytetrade.io/web3os/installer/pkg/version/kubesphere/templates"
	"github.com/pkg/errors"
)

// ~ CopyFiles
type CopyFiles struct {
	common.KubeAction
}

func (t *CopyFiles) Execute(runtime connector.Runtime) error {
	fmt.Println("---1---", runtime.GetRunner().Host.GetName())
	fmt.Println("---2---", runtime.RemoteHost().GetName())
	fmt.Println("---3---", runtime.GetHostWorkDir())
	var cmd = "/usr/local/bin/kubectl get pods -A"
	if _, err := runtime.GetRunner().SudoCmd(cmd, false, true); err != nil {
		return err
	}
	fmt.Println("---4---")

	var kubeConfigFile = path.Join(runtime.GetHostWorkDir(), templates.KsInstaller.Name())
	ksStr, err := util.Render(templates.KsInstaller, util.Data{})
	if err != nil {
		return err
	}

	if err := util.WriteFile(kubeConfigFile, []byte(ksStr), cc.FileMode0644); err != nil {
		return errors.Wrap(errors.WithStack(err), fmt.Sprintf("write kubesphere crd %s failed", kubeConfigFile))
	}

	var clusterConfigFile = path.Join(runtime.GetHostWorkDir(), "cluster-config.yaml")
	configurationBase64 := base64.StdEncoding.EncodeToString([]byte(t.KubeConf.Cluster.KubeSphere.Configurations))

	if err := util.WriteFile(clusterConfigFile, []byte(configurationBase64), cc.FileMode0644); err != nil {
		return errors.Wrap(errors.WithStack(err), fmt.Sprintf("write kubesphere crd %s failed", kubeConfigFile))
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

	copyFiles := &task.RemoteTask{
		Name:   "CopyFiles",
		Hosts:  m.Runtime.GetHostsByRole(common.Master),
		Action: new(CopyFiles),
	}

	m.Tasks = []task.Interface{
		checkMacCommandExists,
		copyFiles,
	}
}
