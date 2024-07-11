package plugins

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"github.com/pkg/errors"
)

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
	return nil
}
