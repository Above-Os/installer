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

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
)

type PackageDownload struct {
	common.KubeAction
}

func (d *PackageDownload) Execute(runtime connector.Runtime) error {
	// todo download
	fmt.Println("---PackageDownload / Execute---")

	var arch = "amd64"
	DownloadPackage(d.KubeConf, runtime.GetWorkDir(), "0.0.1", arch, d.PipelineCache)
	return nil
}
