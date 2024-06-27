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
		Name:   "DownloadPackage",
		Desc:   "Download installer packages",
		Action: new(PackageDownload),
		Retry:  0,
	}

	untar := &task.LocalTask{
		Name:   "UntarPackage",
		Desc:   "Untar installer package",
		Action: new(PackageUntar),
		Retry:  0,
	}

	m.Tasks = []task.Interface{download, untar}
}
