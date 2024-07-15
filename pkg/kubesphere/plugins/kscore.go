package plugins

import (
	"fmt"
	"strconv"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/prepare"
	"bytetrade.io/web3os/installer/pkg/core/task"
	"github.com/pkg/errors"
)

var kscorecrds = []map[string]string{
	{
		"ns":       "kubesphere-controls-system",
		"kind":     "serviceaccounts",
		"resource": "kubesphere-cluster-admin",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-controls-system",
		"kind":     "serviceaccounts",
		"resource": "kubesphere-router-serviceaccount",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-controls-system",
		"kind":     "role",
		"resource": "system:kubesphere-router-role",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-controls-system",
		"kind":     "rolebinding",
		"resource": "nginx-ingress-role-nisa-binding",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-controls-system",
		"kind":     "deployment",
		"resource": "default-http-backend",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-controls-system",
		"kind":     "service",
		"resource": "default-http-backend",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "secrets",
		"resource": "ks-controller-manager-webhook-cert",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "serviceaccounts",
		"resource": "kubesphere",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "configmaps",
		"resource": "ks-router-config",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "configmaps",
		"resource": "sample-bookinfo",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "clusterroles",
		"resource": "system:kubesphere-router-clusterrole",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "clusterrolebindings",
		"resource": "system:nginx-ingress-clusterrole-nisa-binding",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "clusterrolebindings",
		"resource": "system:kubesphere-cluster-admin",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "clusterrolebindings",
		"resource": "kubesphere",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "services",
		"resource": "ks-apiserver",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "services",
		"resource": "ks-controller-manager",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "deployments",
		"resource": "ks-apiserver",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "deployments",
		"resource": "ks-controller-manager",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "validatingwebhookconfigurations",
		"resource": "users.iam.kubesphere.io",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "validatingwebhookconfigurations",
		"resource": "resourcesquotas.quota.kubesphere.io",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "validatingwebhookconfigurations",
		"resource": "network.kubesphere.io",
		"release":  "ks-core",
	},
	{
		"ns":       "kubesphere-system",
		"kind":     "users.iam.kubesphere.io",
		"resource": "admin",
		"release":  "ks-core",
	},
}

// ~ PacthKsCore
type PacthKsCore struct {
	common.KubeAction
}

func (t *PacthKsCore) Execute(runtime connector.Runtime) error {
	var secretsNum int64
	var crdNum int64
	var secretsNumIf, ok = t.PipelineCache.Get(common.CacheSecretsNum)
	if ok && secretsNumIf != nil {
		secretsNum = secretsNumIf.(int64)
	}

	crdNumIf, ok := t.PipelineCache.Get(common.CacheCrdsNUm)
	if ok && crdNumIf != nil {
		crdNum = crdNumIf.(int64)
	}

	var kubectl = "/usr/local/bin/kubectl"
	if secretsNum == 0 && crdNum != 0 {
		for _, item := range kscorecrds {
			var cmd = fmt.Sprintf("%s -n %s annotate --overwrite %s %s meta.helm.sh/release-name=%s && %s -n %s annotate --overwrite %s %s meta.helm.sh/release-namespace=%s && %s -n %s label --overwrite %s %s app.kubernetes.io/managed-by=Helm",
				kubectl, item["ns"], item["kind"], item["resource"], item["release"], kubectl, item["ns"], item["kind"], item["resource"], kubectl,
				item["ns"], item["kind"], item["resource"])

			if _, err := runtime.GetRunner().SudoCmd(cmd, false, true); err != nil {
				return errors.Wrap(errors.WithStack(err), "patch ks-core crd")
			}
		}
	}

	return nil
}

// ~ CheckKsCoreExist
type CheckKsCoreExist struct {
	common.KubeAction
}

func (t *CheckKsCoreExist) Execute(runtime connector.Runtime) error {
	var cmd string

	cmd = fmt.Sprintf("/usr/local/bin/kubectl -n %s get secrets --field-selector=type=helm.sh/release.v1  | grep ks-core |wc -l", common.NamespaceKubesphereSystem)
	stdout, _ := runtime.GetRunner().SudoCmd(cmd, false, false)

	secretNum, err := strconv.ParseInt(stdout, 10, 64)
	if err != nil {
		secretNum = 0
	}

	cmd = "/usr/local/bin/kubectl get crd users.iam.kubesphere.io  | grep 'users.iam.kubesphere.io' |wc -l"
	stdout, _ = runtime.GetRunner().SudoCmd(cmd, false, false)

	usersCrdNum, err := strconv.ParseInt(stdout, 10, 64)
	if err != nil {
		usersCrdNum = 0
	}

	logger.Debugf("secretNum: %d, usersCrdNum: %d", secretNum, usersCrdNum)

	t.ModuleCache.Set(common.CacheSecretsNum, secretNum)
	t.ModuleCache.Set(common.CacheCrdsNUm, usersCrdNum)

	return nil
}

// +

type DeployKsCoreModule struct {
	common.KubeModule
}

func (m *DeployKsCoreModule) Init() {
	m.Name = "DeployKsCore"

	checkKsCoreExist := &task.RemoteTask{
		Name:  "CheckKsCoreExist",
		Hosts: m.Runtime.GetHostsByRole(common.ETCD),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
			new(NotEqualDesiredVersion),
			new(common.GetMasterNum),
		},
		Action:   new(CheckKsCoreExist),
		Parallel: false,
		Retry:    0,
	}

	pacthKsCore := &task.RemoteTask{
		Name:  "PacthKsCore",
		Hosts: m.Runtime.GetHostsByRole(common.ETCD),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
			new(NotEqualDesiredVersion),
			new(common.GetMasterNum),
		},
		Action:   new(PacthKsCore),
		Parallel: false,
		Retry:    0,
	}

	m.Tasks = []task.Interface{
		checkKsCoreExist,
		pacthKsCore,
	}
}
