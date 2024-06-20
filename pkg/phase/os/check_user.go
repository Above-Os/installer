package os

import (
	"bytetrade.io/web3os/installer/pkg/common"
)

func CheckCurrentUserPipeline() error {
	_, err := common.NewLocalRuntime(false, false)
	if err != nil {
		return err
	}

	return nil
}
