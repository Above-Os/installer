package terminus

import (
	"path"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/action"
	corecommon "bytetrade.io/web3os/installer/pkg/core/common"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/task"
	"bytetrade.io/web3os/installer/pkg/scripts"
)

// ~ ClearIPsModule
// called when performing the uninstall operation
type ClearIPsModule struct {
	common.KubeModule
}

func (m *ClearIPsModule) Init() {
	m.Name = "ClearIPs"

	killContainerd := &task.RemoteTask{
		Name:     "KillContainerd",
		Hosts:    m.Runtime.GetHostsByRole(common.Master),
		Action:   new(KillContainerd),
		Parallel: true,
		Retry:    0,
	}

	m.Tasks = []task.Interface{
		killContainerd,
	}
}

// ~ UninstallTerminusCliModule
type UninstallTerminusCliModule struct {
	module.BaseTaskModule
}

func (m *UninstallTerminusCliModule) Init() {
	m.Name = "Uninstall Terminus Cli Mode"

	copyUninstallScript := &task.LocalTask{
		Name:   "Copy",
		Action: new(scripts.CopyUninstallScriptTask),
	}

	checkFileExists := &task.LocalTask{
		Name:   "Check Script Exists",
		Action: new(CheckUninstallScriptExistsAction),
	}

	m.Tasks = []task.Interface{
		copyUninstallScript,
		checkFileExists,
	}
}

// ~ ExecUninstallScriptModule
// todo 测试阶段，执行卸载脚本
type ExecUninstallScriptModule struct {
	module.BaseTaskModule
}

func (m *ExecUninstallScriptModule) Init() {
	m.Name = "ExecUninstallScript"

	var proxy, _ = m.PipelineCache.Get("proxy")
	var kubeType, _ = m.PipelineCache.Get("kube_type")

	var envs = make(map[string]string)
	envs["PROXY"] = proxy.(string)
	envs["KUBE_TYPE"] = kubeType.(string)
	envs["FORCE_UNINSTALL_CLUSTER"] = "1"

	execUninstallScript := &task.LocalTask{
		Name: "ExecUninstallScript",
		Action: &action.Script{
			File:        path.Join(m.Runtime.GetPackageDir(), corecommon.InstallDir, corecommon.UninstallOsScript),
			Args:        []string{},
			Envs:        envs,
			PrintOutput: false,
			PrintLine:   true,
			Ignore:      false,
		},
		Retry: 0,
	}

	m.Tasks = []task.Interface{
		execUninstallScript,
	}
}
