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

package binaries

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/task"
	"github.com/pkg/errors"
)

// ~ PatchUbuntu24AppArmorModule
type PatchUbuntu24AppArmorModule struct {
	module.BaseTaskModule
}

func (m *PatchUbuntu24AppArmorModule) GetName() string {
	return "PatchUbuntu24AppArmorModule"
}

func (m *PatchUbuntu24AppArmorModule) Init() {
	m.Name = "PatchUbuntu24AppArmorModule"
	m.Desc = "Patch Ubuntu 24.04 AppArmor"

	appArmorDownload := &task.LocalTask{
		Name:   "PatchUbuntu24AppArmorModule",
		Desc:   "Setup App Armor for Ubuntu 24.x",
		Action: new(AppArmorDownload),
	}

	appArmorInstall := &task.LocalTask{
		Name:   "PatchUbuntu24AppArmorModule",
		Desc:   "Setup App Armor for Ubuntu 24.x",
		Action: new(AppArmorInstall),
	}

	m.Tasks = []task.Interface{
		appArmorDownload,
		appArmorInstall,
	}
}

// ~ NodeBinariesModule
type NodeBinariesModule struct {
	common.KubeModule
}

func (n *NodeBinariesModule) GetName() string {
	return "NodeBinariesModule"
}

func (n *NodeBinariesModule) Init() {
	n.Name = "NodeBinariesModule"
	n.Desc = "Download installation binaries"

	download := &task.LocalTask{
		Name:   "DownloadBinaries",
		Desc:   "Download installation binaries",
		Action: new(Download),
	}

	n.Tasks = []task.Interface{
		download,
	}
}

type K3sNodeBinariesModule struct {
	common.KubeModule
}

func (k *K3sNodeBinariesModule) GetName() string {
	return "K3sNodeBinariesModule"
}

func (k *K3sNodeBinariesModule) Init() {
	k.Name = "K3sNodeBinariesModule"
	k.Desc = "Download installation binaries"

	download := &task.LocalTask{
		Name:   "DownloadBinaries",
		Desc:   "Download installation binaries",
		Action: new(K3sDownload),
	}

	k.Tasks = []task.Interface{
		download,
	}
}

type ArtifactBinariesModule struct {
	common.ArtifactModule
}

func (a *ArtifactBinariesModule) GetName() string {
	return "ArtifactBinariesModule"
}

func (a *ArtifactBinariesModule) Init() {
	a.Name = "ArtifactBinariesModule"
	a.Desc = "Download artifact binaries"

	download := &task.LocalTask{
		Name:   "DownloadBinaries",
		Desc:   "Download manifest expect binaries",
		Action: new(ArtifactDownload),
	}

	a.Tasks = []task.Interface{
		download,
	}
}

type K3sArtifactBinariesModule struct {
	common.ArtifactModule
}

func (a *K3sArtifactBinariesModule) Init() {
	a.Name = "K3sArtifactBinariesModule"
	a.Desc = "Download artifact binaries"

	download := &task.LocalTask{
		Name:   "K3sDownloadBinaries",
		Desc:   "Download k3s manifest expect binaries",
		Action: new(K3sArtifactDownload),
	}

	a.Tasks = []task.Interface{
		download,
	}
}

type RegistryPackageModule struct {
	common.KubeModule
}

func (n *RegistryPackageModule) GetName() string {
	return "RegistryPackageModule"
}

func (n *RegistryPackageModule) Init() {
	n.Name = "RegistryPackageModule"
	n.Desc = "Download registry package"

	if len(n.Runtime.GetHostsByRole(common.Registry)) == 0 {
		logger.Fatal(errors.New("[registry] node not found in the roleGroups of the configuration file"))
	}

	download := &task.LocalTask{
		Name:   "DownloadRegistryPackage",
		Desc:   "Download registry package",
		Action: new(RegistryPackageDownload),
	}

	n.Tasks = []task.Interface{
		download,
	}
}

type CriBinariesModule struct {
	common.KubeModule
}

func (i *CriBinariesModule) GetName() string {
	return "CriBinariesModule"
}

func (i *CriBinariesModule) Init() {
	i.Name = "CriBinariesModule"
	i.Desc = "Download Cri package"
	switch i.KubeConf.Arg.Type {
	case common.Docker:
		i.Tasks = CriBinaries(i)
	case common.Conatinerd:
		i.Tasks = CriBinaries(i)
	default:
	}

}

func CriBinaries(p *CriBinariesModule) []task.Interface {

	download := &task.LocalTask{
		Name:   "DownloadCriPackage",
		Desc:   "Download Cri package",
		Action: new(CriDownload),
	}

	p.Tasks = []task.Interface{
		download,
	}
	return p.Tasks
}
