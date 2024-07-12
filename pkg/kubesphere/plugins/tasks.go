package plugins

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/utils"
	"github.com/pkg/errors"
)

// ~ CopyEmbedFiles
type CopyEmbedFiles struct {
	common.KubeAction
}

func (t *CopyEmbedFiles) Execute(runtime connector.Runtime) error {
	return utils.CopyEmbed(Assets(), "files", runtime.GetFilesDir())
}

// ! moved to prepares
// // ~ CheckMasterNum
// type CheckMasterNum struct {
// 	common.KubeAction
// }

// func (t *CheckMasterNum) Execute(runtime connector.Runtime) error {
// 	var cmd = fmt.Sprintf("/usr/local/bin/kubectl get node | awk '{if(NR>1){print $3}}' | grep master | wc -l")
// 	var stdout, err = runtime.GetRunner().SudoCmd(cmd, false, true)
// 	if err != nil {
// 		return errors.Wrap(errors.WithStack(err), "get master num failed")
// 	}

// 	var enableHA = "0"
// 	if stdout != "" && stdout != "0" && stdout != "1" {
// 		enableHA = "1"
// 	}
// 	t.PipelineCache.Set(common.CacheEnableHA, enableHA)
// 	return nil
// }

// ~ InitNamespace
type InitNamespace struct {
	common.KubeAction
}

func (t *InitNamespace) Execute(runtime connector.Runtime) error {
	_, err := runtime.GetRunner().SudoCmd(
		fmt.Sprintf(`cat <<EOF | /usr/local/bin/kubectl apply -f -
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
EOF
`, common.NamespaceKubesphereControlsSystem, common.NamespaceKubesphereMonitoringFederated), false, true)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), fmt.Sprintf("create namespace: %s and %s",
			common.NamespaceKubesphereControlsSystem, common.NamespaceKubesphereMonitoringFederated))
	}

	var allNs = []string{
		common.NamespaceDefault,
		common.NamespaceKubeNodeLease,
		common.NamespaceKubePublic,
		common.NamespaceKubeSystem,
		common.NamespaceKubekeySystem,
		common.NamespaceKubesphereControlsSystem,
		common.NamespaceKubesphereMonitoringFederated,
		common.NamespaceKubesphereMonitoringSystem,
		common.NamespaceKubesphereSystem,
	}

	for _, ns := range allNs {
		if _, err := runtime.GetRunner().SudoCmd(fmt.Sprintf("/usr/local/bin/kubectl label ns %s kubesphere.io/workspace=system-workspace --overwrite", ns), false, true); err != nil {
			logger.Errorf("label ns %s kubesphere.io/workspace=system-workspace failed: %v", ns, err)
			return errors.Wrap(errors.WithStack(err), fmt.Sprintf("label namespace %s kubesphere.io/workspace=system-workspace failed: %v", ns, err))
		}

		if _, err := runtime.GetRunner().SudoCmd(fmt.Sprintf("/usr/local/bin/kubectl label ns %s kubesphere.io/namespace=%s --overwrite", ns, ns), false, true); err != nil {
			logger.Errorf("label ns %s kubesphere.io/namespace=%s failed: %v", ns, ns, err)
			return errors.Wrap(errors.WithStack(err), fmt.Sprintf("label namespace %s kubesphere.io/namespace=%s failed: %v", ns, ns, err))
		}
	}

	return nil
}
