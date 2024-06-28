package terminus

import (
	"fmt"
	"path"

	"bytetrade.io/web3os/installer/pkg/core/action"
	"bytetrade.io/web3os/installer/pkg/core/common"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/task"
)

// ~ UninstallTerminusCliModule
type UninstallTerminusCliModule struct {
	module.BaseTaskModule
}

func (m *UninstallTerminusCliModule) Init() {
	m.Name = "Uninstall Terminus Cli Mode"

	aa, bb := m.PipelineCache.Get("hello")
	fmt.Printf("aa: %s, bb: %v\n", aa, bb)

	checkFileExists := &task.LocalTask{
		Name:    "Check Script Exists",
		Prepare: new(CheckUninstallScriptExists), // todo 这里缺少actioin 会报错
	}

	m.Tasks = []task.Interface{
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

	a, b := m.ModuleCache.Get("uninstall")
	fmt.Printf("---1--- %+v  %v\n", a, b)

	execUninstallScript := &task.LocalTask{
		Name: "ExecUninstallScript",
		Action: &action.Script{
			Name:        "ExecUninstallScript",
			File:        path.Join(m.Runtime.GetPackageDir(), common.InstallDir, common.UninstallOsScript),
			Args:        []string{},
			Envs:        make(map[string]string),
			PrintOutput: true,
			PrintLine:   true,
			Ignore:      true,
		},
		Retry: 0,
	}

	m.Tasks = []task.Interface{
		execUninstallScript,
	}
}
