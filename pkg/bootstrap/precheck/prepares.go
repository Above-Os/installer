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

package precheck

import (
	"fmt"
	"net"
	"os/exec"
	"strings"

	kubekeyapiv1alpha2 "bytetrade.io/web3os/installer/apis/kubekey/v1alpha2"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/prepare"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"bytetrade.io/web3os/installer/pkg/files"
	"github.com/pkg/errors"
)

// ~ DownloadDepsExt
type DownloadDepsExt struct {
	prepare.BasePrepare
}

func (p *DownloadDepsExt) PreCheck(runtime connector.Runtime) (bool, error) {
	logger.Debug("[prepare] DownloadDepsExt")
	binaries := []*files.KubeBinary{}

	// switch constants.OsPlatform {
	// case common.Ubuntu:
	// 	if strings.HasPrefix(constants.OsVersion, "24.") {
	// 		apparmor := files.NewKubeBinary("apparmor", kubekeyapiv1alpha2.DefaultArch, kubekeyapiv1alpha2.DefaultUbuntu24AppArmonVersion, runtime.GetDependDir())
	// 		binaries = append(binaries, apparmor)
	// 	}
	// case common.Debian, common.Raspbian, common.CentOs, common.Fedora, common.RHEl:
	// default:
	// 	socat := files.NewKubeBinary("socat", kubekeyapiv1alpha2.DefaultArch, kubekeyapiv1alpha2.DefaultSocatVersion, runtime.GetDependDir())
	// 	contrack := files.NewKubeBinary("contrack", kubekeyapiv1alpha2.DefaultArch, kubekeyapiv1alpha2.DefaultConntrackVersion, runtime.GetDependDir())
	// 	binaries = append(binaries, socat, contrack)
	// }

	apparmor := files.NewKubeBinary("apparmor", kubekeyapiv1alpha2.DefaultArch, kubekeyapiv1alpha2.DefaultUbuntu24AppArmonVersion, runtime.GetDependDir())
	binaries = append(binaries, apparmor)
	socat := files.NewKubeBinary("socat", kubekeyapiv1alpha2.DefaultArch, kubekeyapiv1alpha2.DefaultSocatVersion, runtime.GetDependDir())
	contrack := files.NewKubeBinary("conntrack", kubekeyapiv1alpha2.DefaultArch, kubekeyapiv1alpha2.DefaultConntrackVersion, runtime.GetDependDir())
	binaries = append(binaries, socat, contrack)

	binariesMap := make(map[string]*files.KubeBinary)

	for _, binary := range binaries {
		if err := binary.CreateBaseDir(); err != nil {
			return false, errors.Wrapf(errors.WithStack(err), "create file %s base dir failed", binary.FileName)
		}

		logger.Infof("%s downloading %s %s %s ...", common.LocalHost, constants.OsArch, binary.ID, binary.Version)

		binariesMap[binary.ID] = binary
		if util.IsExist(binary.Path()) {
			p := binary.Path()
			if err := binary.SHA256Check(); err != nil {
				_ = exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f %s", p)).Run()
			} else {
				continue
			}
		}

		if err := binary.Download(); err != nil {
			return false, fmt.Errorf("Failed to download %s binary: %s error: %w ", binary.ID, binary.Url, err)
		}
	}

	p.PipelineCache.Set(common.KubeBinaries+"-"+constants.OsArch, binariesMap)

	return true, nil
}

// ~ LocalIpCheck
type LocalIpCheck struct {
	prepare.BasePrepare
}

func (p *LocalIpCheck) PreCheck(runtime connector.Runtime) (bool, error) {
	var localIp = constants.LocalIp
	ip := net.ParseIP(localIp)
	if ip == nil {
		return false, fmt.Errorf("invalid local ip %s", localIp)
	}

	if ip4 := ip.To4(); ip4 == nil {
		return false, fmt.Errorf("invalid local ip %s", localIp)
	}

	switch localIp {
	case "172.17.0.1", "127.0.0.1", "127.0.1.1":
		return false, fmt.Errorf("invalid local ip %s", localIp)
	default:
	}
	return true, nil
}

// ~ OsSupportCheck
type OsSupportCheck struct {
	prepare.BasePrepare
}

func (p *OsSupportCheck) PreCheck(runtime connector.Runtime) (bool, error) {
	switch constants.OsType {
	case common.Linux:
		switch constants.OsPlatform {
		case common.Ubuntu:
			if strings.HasPrefix(constants.OsVersion, "20.") || strings.HasPrefix(constants.OsVersion, "22.") || strings.HasPrefix(constants.OsVersion, "24.") {
				return true, nil
			}
			return false, fmt.Errorf("os %s version %s not support", constants.OsPlatform, constants.OsVersion)
		case common.Debian:
			if strings.HasPrefix(constants.OsVersion, "11") || strings.HasPrefix(constants.OsVersion, "12") {
				return true, nil
			}
			return false, fmt.Errorf("os %s version %s not support", constants.OsPlatform, constants.OsVersion)
		default:
			return false, fmt.Errorf("platform %s not support", constants.OsPlatform)
		}
	default:
		return false, fmt.Errorf("os %s not support", constants.OsType)
	}
}

// ~ KubeSphereExist
type KubeSphereExist struct {
	common.KubePrepare
}

func (k *KubeSphereExist) PreCheck(runtime connector.Runtime) (bool, error) {
	currentKsVersion, ok := k.PipelineCache.GetMustString(common.KubeSphereVersion)
	if !ok {
		return false, errors.New("get current KubeSphere version failed by pipeline cache")
	}
	if currentKsVersion != "" {
		return true, nil
	}
	return false, nil
}
