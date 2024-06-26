/*
 Copyright 2021 The KubeSphere Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package packages

import (
	"fmt"

	kubekeyapiv1alpha2 "bytetrade.io/web3os/installer/apis/kubekey/v1alpha2"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/util"
)

type PackageDownload struct {
	common.KubeAction
}

func (d *PackageDownload) Execute(runtime connector.Runtime) error {
	logger.Debug("[action] PackageDownload")

	if err := DownloadInstallPackage(d.KubeConf, runtime.GetPackageDir(), "0.0.1", kubekeyapiv1alpha2.DefaultArch, d.PipelineCache); err != nil {
		return err
	}
	return nil
}

// todo 这里是一个测试解压 full 包的 action
type PackageUntar struct {
	common.KubeAction
}

func (a *PackageUntar) Execute(runtime connector.Runtime) error {
	var pkgFile = fmt.Sprintf("%s/install-wizard-full_%s.tar.gz", runtime.GetPackageDir(), kubekeyapiv1alpha2.DefaultArch)
	if ok := util.IsExist(pkgFile); !ok {
		return fmt.Errorf("package %s not exist", pkgFile)
	}

	if err := util.RemoveDir("/home/zhaoyu/install-wizard"); err != nil {
		return fmt.Errorf("remove %s failed %v", "/home/zhaoyu/install-wizard", err)
	}

	if err := util.Untar(pkgFile, "/home/zhaoyu/install-wizard"); err != nil {
		return fmt.Errorf("untar %s failed %v", pkgFile, err)
	}
	logger.Infof("untar %s success", pkgFile)
	return nil
}
