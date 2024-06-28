package cli

import "bytetrade.io/web3os/installer/pkg/common"

func Uninstall() error {
	runtime, err := common.NewLocalRuntime(false, false)
	if err != nil {
		return err
	}

	runtime.GetRunner().Host.Echo()

	return nil
}

func Install() {

}

func Restore() {

}
