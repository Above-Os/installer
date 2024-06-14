package packages

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/task"
)

type PackagesModule struct {
	common.KubeModule
}

func (m *PackagesModule) Init() {
	m.Name = "PackageModule"
	m.Desc = "Download installer packages"

	download := &task.LocalTask{
		Name:   "DownloadPackages",
		Desc:   "Download installer packages",
		Action: new(PackageDownload),
	}

	m.Tasks = []task.Interface{download}
}
