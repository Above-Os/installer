package options

import (
	"fmt"
	"os"

	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/utils"
)

func InitEnv(o *ApiOptions) {
	fmt.Println(constants.Logo)

	workDir, err := utils.WorkDir()
	if err != nil {
		fmt.Println("working path error",	err)
		os.Exit(1)
	}

	constants.WorkDir = workDir
	constants.ApiServerListenAddress = o.Port
	constants.Proxy = o.Proxy
}