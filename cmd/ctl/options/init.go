package options

import (
	"fmt"
	"os"
	"runtime"

	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"bytetrade.io/web3os/installer/pkg/utils"
)

func InitEnv(o *ApiOptions) {
	fmt.Println(constants.Logo)

	workDir, err := utils.WorkDir()
	if err != nil {
		fmt.Println("working path error", err)
		os.Exit(1)
	}

	constants.WorkDir = workDir
	constants.ApiServerListenAddress = o.Port
	constants.Proxy = o.Proxy
}

func InitLocal() {
	constants.LocalIp = util.LocalIP()
	constants.OsArch = runtime.GOARCH
}
