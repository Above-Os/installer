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

// ~ WizardTask
type WizardTask struct {
	common.KubeAction
}

// todo 这里增加一个检查 terminus 是否安装完成的检查，需要重启 wizard 这容器
func (t *WizardTask) Execute(runtime connector.Runtime) error {
	return nil
}


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

	var provider = runtime.GetStorage()
	var domainName = installReq.Config.DomainName
	var userName = installReq.Config.UserName
	var kubeType = strings.ToLower(installReq.Config.KubeType)
	var d = path.Join(runtime.GetPackageDir(), corecommon.InstallDir)
	var installCommand = fmt.Sprintf("export TERMINUS_OS_DOMAINNAME=%s;export TERMINUS_OS_USERNAME=%s;export KUBE_TYPE=%s;bash %s/install_cmd.sh", domainName, userName, kubeType, d)

	var out = make(chan []interface{}, 10)
	go func() {
		_, _, err = util.ExecWithChannel(installCommand, false, true, out)
		if err != nil {
			return
		}
	}()

	for {
		select {
		case r, ok := <-out:
			if !ok {
				break
			}
			if r == nil || len(r) != 2 {
				continue
			}

			msg := r[0].(string)
			percent := r[1].(int64)
			state := corecommon.StateInstall

			pos := strings.Index(msg, "[INFO]")
			if pos >= 0 {
				msg = msg[pos+len("[INFO]"):]
			} else {
				msg = strings.ReplaceAll(msg, "[INFO]", "")
			}
			msg = strings.ReplaceAll(msg, "\u001b", "")
			msg = strings.ReplaceAll(msg, "[0m", "")
			msg = strings.TrimSpace(msg)
			if percent == corecommon.DefaultInstallSteps {
				state = corecommon.StateSuccess
			}
			// fmt.Printf("---1--- [%s]  [%d]\n", msg, percent)
			if strings.Contains(msg, "installing k8s and kubesphere") {
				msg = "installing k8s and kubesphere, this will take a few minutes, please wait ..."
			}
			if err := provider.SaveInstallLog(msg, state, int64(percent*10000/corecommon.DefaultInstallSteps)); err != nil {
				logger.Errorf("save install log failed %v", err)
			}
		}
	}

	return err
}
