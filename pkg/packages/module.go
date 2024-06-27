package packages

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/task"
	"bytetrade.io/web3os/installer/pkg/model"
)

type PackagesModule struct {
	common.KubeModule
	common.KubeConf
}

func (m *PackagesModule) Init() {
	m.Name = "PackageModule"
	m.Desc = "Download installer packages"
	var installReq model.InstallModelReq
	var ok bool

	tmp := m.KubeConf.Arg.Request
	if installReq, ok = any(tmp).(model.InstallModelReq); !ok {
		return
	}

	fmt.Printf("---packagesModule--- %+v\n", installReq)

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
