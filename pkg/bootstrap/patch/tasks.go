package patch

import (
	"fmt"
	"os/exec"
	"path"

	kubekeyapiv1alpha2 "bytetrade.io/web3os/installer/apis/kubekey/v1alpha2"
	"bytetrade.io/web3os/installer/pkg/binaries"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/core/action"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"bytetrade.io/web3os/installer/pkg/files"
	"github.com/pkg/errors"
)

// ~ PatchTask apt-get install
type PatchTask struct {
	action.BaseAction
}

func (t *PatchTask) Execute(runtime connector.Runtime) error {
	var host = runtime.GetRunner().Host
	var cmd string
	var debianFrontend string
	var pre_reqs = "apt-transport-https ca-certificates curl"

	if _, err := runtime.GetRunner().Host.GetCommand(common.CommandGPG); err != nil {
		pre_reqs = pre_reqs + " gnupg"
	}

	switch constants.OsPlatform {
	case common.Debian:
		debianFrontend = "DEBIAN_FRONTEND=noninteractive"
		fallthrough
	case common.Ubuntu, common.Raspbian:
		if _, _, err := host.Exec(fmt.Sprintf("%s update -qq", constants.PkgManager), false, false); err != nil {
			logger.Errorf("update os error %v", err)
			return err
		}

		if _, _, err := host.Exec("apt --fix-broken install -y", false, true); err != nil {
			logger.Errorf("fix-broken install error %v", err)
			return err
		}

		if _, _, err := host.Exec(fmt.Sprintf("%s %s install -y -qq %s", debianFrontend, constants.PkgManager, pre_reqs), false, true); err != nil {
			logger.Errorf("install deps %s error %v", pre_reqs, err)
			return err
		}

		var cmd = "conntrack socat apache2-utils ntpdate net-tools make gcc openssh-server bison flex"
		if _, _, err := host.Exec(fmt.Sprintf("%s %s install -y %s", debianFrontend, constants.PkgManager, cmd), false, true); err != nil {
			logger.Errorf("install deps %s error %v", cmd, err)
			return err
		}
	case common.CentOs, common.Fedora, common.RHEl:
		cmd = "conntrack socat httpd-tools ntpdate net-tools make gcc openssh-server"
		if _, _, err := host.Exec(fmt.Sprintf("%s install -y %s", constants.PkgManager, cmd), false, true); err != nil {
			logger.Errorf("install deps %s error %v", cmd, err)
			return err
		}
	}

	return nil
}

// ~ SocatTask
type SocatTask struct {
	action.BaseAction
}

func (t *SocatTask) Execute(runtime connector.Runtime) error {
	filePath, fileName, err := binaries.DownloadSocat(runtime.GetWorkDir(), kubekeyapiv1alpha2.DefaultSocatVersion, constants.OsArch, t.PipelineCache)
	if err != nil {
		logger.Errorf("failed to download socat: %v", err)
		return err
	}
	f := path.Join(filePath, fileName)
	if _, _, err := runtime.GetRunner().Host.Exec(fmt.Sprintf("tar xzvf %s -C %s", f, filePath), false, false); err != nil {
		logger.Errorf("failed to extract %s %v", f, err)
		return err
	}

	tp := path.Join(filePath, fmt.Sprintf("socat-%s", kubekeyapiv1alpha2.DefaultSocatVersion))
	if err := runtime.GetRunner().Host.ChangeDir(tp); err == nil {
		if _, _, err := runtime.GetRunner().Host.Exec("./configure --prefix=/usr && make -j4 && make install && strip socat", false, false); err != nil {
			logger.Errorf("failed to install socat %v", err)
			return err
		}
	}
	if err := runtime.GetRunner().Host.ChangeDir(runtime.GetRootDir()); err != nil {
		logger.Errorf("failed to change dir %v", err)
		return err
	}

	return nil
}

// ~ ConntrackTask
type ConntrackTask struct {
	action.BaseAction
}

func (t *ConntrackTask) Execute(runtime connector.Runtime) error {
	flexFilePath, flexFileName, err := binaries.DownloadFlex(runtime.GetWorkDir(), kubekeyapiv1alpha2.DefaultFlexVersion, constants.OsArch, t.PipelineCache)
	if err != nil {
		logger.Errorf("failed to download flex: %v", err)
		return err
	}
	filePath, fileName, err := binaries.DownloadConntrack(runtime.GetWorkDir(), kubekeyapiv1alpha2.DefaultConntrackVersion, constants.OsArch, t.PipelineCache)
	if err != nil {
		logger.Errorf("failed to download conntrack: %v", err)
		return err
	}
	fl := path.Join(flexFilePath, flexFileName)
	f := path.Join(filePath, fileName)

	if _, _, err := runtime.GetRunner().Host.Exec(fmt.Sprintf("tar xzvf %s -C %s", fl, filePath), false, true); err != nil {
		logger.Errorf("failed to extract %s %v", flexFilePath, err)
		return err
	}

	if _, _, err := runtime.GetRunner().Host.Exec(fmt.Sprintf("tar xzvf %s -C %s", f, filePath), false, true); err != nil {
		logger.Errorf("failed to extract %s %v", f, err)
		return err
	}

	// install
	fp := path.Join(flexFilePath, fmt.Sprintf("flex-%s", kubekeyapiv1alpha2.DefaultFlexVersion))
	if err := runtime.GetRunner().Host.ChangeDir(fp); err == nil {
		if _, _, err := runtime.GetRunner().Host.Exec("autoreconf -i && ./configure --prefix=/usr && make -j4 && make install", false, true); err != nil {
			logger.Errorf("failed to install flex %v", err)
			return err
		}
	}

	tp := path.Join(filePath, fmt.Sprintf("conntrack-tools-conntrack-tools-%s", kubekeyapiv1alpha2.DefaultConntrackVersion))
	if err := runtime.GetRunner().Host.ChangeDir(tp); err == nil {
		if _, _, err := runtime.GetRunner().Host.Exec("autoreconf -i && ./configure --prefix=/usr && make -j4 && make install", false, true); err != nil {
			logger.Errorf("failed to install conntrack %v", err)
			return err
		}
	}
	if err := runtime.GetRunner().Host.ChangeDir(runtime.GetRootDir()); err != nil {
		logger.Errorf("failed to change dir %v", err)
		return err
	}

	return nil
}

// ~ PatchDeps
// install socat and contrack on other systemc
type PatchDeps struct {
	action.BaseAction
}

func (t *PatchDeps) Execute(runtime connector.Runtime) error {
	// 如果是特殊的系统，需要通过源代码来安装 socat 和 contrack
	switch constants.OsPlatform {
	case common.Ubuntu, common.Debian, common.Raspbian, common.CentOs, common.Fedora, common.RHEl:
		return nil
	}

	socat := files.NewKubeBinary("socat", constants.OsArch, kubekeyapiv1alpha2.DefaultSocatVersion, runtime.GetDependDir())
	conntrack := files.NewKubeBinary("conntrack", constants.OsArch, kubekeyapiv1alpha2.DefaultConntrackVersion, runtime.GetDependDir())

	binaries := []*files.KubeBinary{socat, conntrack}
	binariesMap := make(map[string]*files.KubeBinary)

	for _, binary := range binaries {
		if err := binary.CreateBaseDir(); err != nil {
			return errors.Wrapf(errors.WithStack(err), "create file %s base dir failed", binary.FileName)
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
			return fmt.Errorf("Failed to download %s binary: %s error: %w ", binary.ID, binary.Url, err)
		}
	}

	t.PipelineCache.Set(common.KubeBinaries+"-"+constants.OsArch, binariesMap)
	return nil
}
