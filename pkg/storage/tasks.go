package storage

import (
	"context"
	"fmt"
	"time"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/model"
)

// ~ task SaveInstallConfigTask
type SaveInstallConfigTask struct {
	common.KubeAction
}

func (t *SaveInstallConfigTask) GetName() string {
	return "SaveInstallConfigTask"
}

func (t *SaveInstallConfigTask) Execute(runtime connector.Runtime) error {
	var installReq model.InstallModelReq
	var ok bool
	if installReq, ok = any(t.KubeConf.Arg.Request).(model.InstallModelReq); !ok {
		return fmt.Errorf("invalid install model req %+v", t.KubeConf.Arg.Request)
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if installReq.DebugSaveConfig == 1 {
		return t.KubeConf.Arg.Provider.SaveInstallConfig(ctx, installReq)
	}
	return nil
}
