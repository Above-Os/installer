package terminus

import (
	"fmt"
	"path"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/action"
	corecommon "bytetrade.io/web3os/installer/pkg/core/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/util"
)

// ~ KillContainerd
type KillContainerd struct {
	common.KubeAction
}

func (t *KillContainerd) Execute(runtime connector.Runtime) error {
	runtime.GetRunner().SudoCmdExt("killall /usr/local/bin/containerd", false, false)
	return nil
}

// ~ CheckUninstallScriptExistsAction
type CheckUninstallScriptExistsAction struct {
	action.BaseAction
}

func (t *CheckUninstallScriptExistsAction) Execute(runtime connector.Runtime) error {
	var fileName = path.Join(runtime.GetPackageDir(), corecommon.InstallDir, corecommon.UninstallOsScript)
	if ok := util.IsExist(fileName); !ok {
		return fmt.Errorf("file %s not exists", fileName)
	}

	return nil
}
