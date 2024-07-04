package patch

import (
	"fmt"
	"os/exec"
	"strings"

	kubekeyapiv1alpha2 "bytetrade.io/web3os/installer/apis/kubekey/v1alpha2"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/core/action"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"bytetrade.io/web3os/installer/pkg/files"
	"github.com/pkg/errors"
)

// ~ PatchDepsTask apt-get install
type PatchDepsTask struct {
	action.BaseAction
	Deps []string
}

func (t *PatchDepsTask) Execute(runtime connector.Runtime) error {
	if t.Deps == nil || len(t.Deps) == 0 {
		return nil
	}
	var apps = strings.Join(t.Deps, " ")
	var cmd = fmt.Sprintf("apt-get install -y %s", apps)
	_, _, err := runtime.GetRunner().Host.Exec(cmd, true, true)
	if err != nil {
		logger.Errorf("install %s error %v", apps, err)
		return err
	}
	return nil
}

// ~ SocatTask
type SocatTask struct {
	action.BaseAction
}

func (t *SocatTask) Execute(runtime connector.Runtime) error {
	return nil
}

// ~ ContrackTask
type ContrackTask struct {
	action.BaseAction
}

func (t *ContrackTask) Execute(runtime connector.Runtime) error {
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
	contrack := files.NewKubeBinary("contrack", constants.OsArch, kubekeyapiv1alpha2.DefaultContrackVersion, runtime.GetDependDir())

	binaries := []*files.KubeBinary{socat, contrack}
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
