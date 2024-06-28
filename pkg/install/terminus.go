package install

import (
	"fmt"
	"path"
	"strings"

	"bytetrade.io/web3os/installer/pkg/common"
	corecommon "bytetrade.io/web3os/installer/pkg/core/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"bytetrade.io/web3os/installer/pkg/model"
)

// ~ Terminus（第一阶段的测试 full 安装）
type Terminus struct {
	common.KubeAction
}

func (a *Terminus) Execute(runtime connector.Runtime) error {
	var installReq model.InstallModelReq
	var err error
	var ok bool
	if installReq, ok = any(a.KubeConf.Arg.Request).(model.InstallModelReq); !ok {
		logger.Errorf("invalid install model req %+v", a.KubeConf.Arg.Request)
		return nil
	}
	if installReq.DebugInstall == 1 {
		var domainName = installReq.DomainName
		var userName = installReq.UserName
		var kubeType = strings.ToLower(installReq.KubeType)
		var d = path.Join(runtime.GetPackageDir(), corecommon.InstallDir)
		var installCommand = fmt.Sprintf("export TERMINUS_OS_DOMAINNAME=%s;export TERMINUS_OS_USERNAME=%s;export KUBE_TYPE=%s;bash %s/install_cmd.sh", domainName, userName, kubeType, d)

		var out chan string = make(chan string)
		go func() {
			_, _, err = util.ExecWithChannel(installCommand, false, true, out)
			if err != nil {
				return
			}
		}()

		for {
			select {
			case _, ok := <-out:
				if !ok {
					break
				}
			}
		}

		// _, _, err := util.Exec(installCommand, true)
		// if err != nil {
		// 	return err
		// }
	}

	return err
}
