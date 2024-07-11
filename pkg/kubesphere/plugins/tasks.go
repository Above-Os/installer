package plugins

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"github.com/pkg/errors"
)

// ~ InitNamespace
type InitNamespace struct {
	common.KubeAction
}

func (t *InitNamespace) Execute(runtime connector.Runtime) error {
	_, err := runtime.GetRunner().SudoCmd(`cat <<EOF | /usr/local/bin/kubectl apply -f -
apiVersion: v1
kind: Namespace
metadata:
  name: kubesphere-controls-system
---
apiVersion: v1
kind: Namespace
metadata:
  name: kubesphere-monitoring-federated
EOF
`, false, true)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), "create namespace: kubesphere-controls-system and kubesphere-monitoring-federated")
	}

	var allNs = []string{
		"default",
		"kube-node-lease",
		"kube-public",
		"kube-system",
		"kubekey-system",
		"kubesphere-controls-system",
		"kubesphere-monitoring-federated",
		"kubesphere-monitoring-system",
		"kubesphere-system",
	}

	for _, ns := range allNs {
		if _, err := runtime.GetRunner().SudoCmd(fmt.Sprintf("/usr/local/bin/kubectl label ns %s kubesphere.io/workspace=system-workspace", ns), false, true); err != nil {
			logger.Errorf("label ns %s kubesphere.io/workspace=system-workspace failed: %v", ns, err)
			return errors.Wrap(errors.WithStack(err), fmt.Sprintf("label namespace %s kubesphere.io/workspace=system-workspace failed: %v", ns, err))
		}

		if _, err := runtime.GetRunner().SudoCmd(fmt.Sprintf("/usr/local/bin/kubectl label ns %s kubesphere.io/namespace=%s", ns, ns), false, true); err != nil {
			logger.Errorf("label ns %s kubesphere.io/namespace=%s failed: %v", ns, ns, err)
			return errors.Wrap(errors.WithStack(err), fmt.Sprintf("label namespace %s kubesphere.io/namespace=%s failed: %v", ns, ns, err))
		}
	}

	return nil
}

// ~ DeploySnapshotController
type DeploySnapshotController struct {
	common.KubeAction
}

func (t *DeploySnapshotController) Execute(runtime connector.Runtime) error {
	fmt.Println("---a---")
	return nil
}
