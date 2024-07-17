package install

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/task"
)

// 这里是测试安装 kk 的
type InstallModule struct {
	common.KubeModule
}

func (m *InstallModule) Init() {
	m.Name = "InstallModule"
	m.Desc = "Install Module"

	// todo 安装这一步要不拆分成多个 action
	checkFileExists := &task.LocalTask{
		Name:   "CheckFileExists",
		Desc:   "check kk exists",
		Action: new(CheckFilesExists),
	}

	copyInstallPackage := &task.LocalTask{
		Name:   "CopyInstallPackage",
		Desc:   "copy install package",
		Action: new(CopyInstallPackage),
	}

	m.Tasks = []task.Interface{checkFileExists, copyInstallPackage}
}
