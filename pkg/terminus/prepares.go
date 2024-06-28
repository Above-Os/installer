package terminus

import (
	"fmt"
	"path"

	"bytetrade.io/web3os/installer/pkg/common"
	corecommon "bytetrade.io/web3os/installer/pkg/core/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/util"
)

type CheckUninstallScriptExists struct {
	common.KubePrepare
}

// 检查卸载脚本是否存在
func (f *CheckUninstallScriptExists) PreCheck(runtime connector.Runtime) (bool, error) {
	var fileName = path.Join(runtime.GetPackageDir(), corecommon.InstallDir, corecommon.UninstallOsScript)
	if ok := util.IsExist(fileName); !ok {
		return false, fmt.Errorf("file %s not exists", fileName)
	}
	return true, nil
}
