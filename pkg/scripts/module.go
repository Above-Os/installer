package scripts

import (
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/task"
)

// ~ CopyUninstallScriptModule
// ! 测试的，目前在做卸载时，原始的 uninstall_cmd.sh 执行会报错，主要还是执行路径的问题
// ! 这里先拷贝内部嵌入的一个修复版本，等后面拆分脚本时再完善
type CopyUninstallScriptModule struct {
	module.BaseTaskModule
}

func (m *CopyUninstallScriptModule) GetName() string {
	return "CopyUninstallScript"
}

func (m *CopyUninstallScriptModule) Init() {
	m.Name = "CopyUninstallScript"
	m.Desc = "copy uninstall script"

	copyUninstallScript := &task.LocalTask{
		Name:   "CopyUninstallScript",
		Desc:   "CopyUninstallScript",
		Action: new(CopyUninstallScriptTask),
	}

	m.Tasks = []task.Interface{
		copyUninstallScript,
	}
}

// ~ CopyScriptsModule
type CopyScriptsModule struct {
	module.BaseTaskModule
}

func (m *CopyScriptsModule) GetName() string {
	return "CopyScriptsModule"
}

func (m *CopyScriptsModule) Init() {
	m.Name = "CopyScriptsModule"
	m.Desc = "Copy scripts"

	copyScripts := &task.LocalTask{
		Name:   "CopyScripts",
		Desc:   "Copy scripts",
		Action: &Copy{},
	}

	greeting := &task.LocalTask{
		Name:        "Greeting",
		Desc:        "Greeting",
		Action:      &Greeting{},
		Retry:       0,
		IgnoreError: true,
	}

	m.Tasks = []task.Interface{
		copyScripts,
		greeting,
	}
}
