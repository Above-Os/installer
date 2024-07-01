package terminus

import (
	"fmt"
	"path"

	"bytetrade.io/web3os/installer/pkg/core/action"
	corecommon "bytetrade.io/web3os/installer/pkg/core/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/util"
)

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
