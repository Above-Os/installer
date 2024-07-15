package plugins

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/prepare"
	"bytetrade.io/web3os/installer/pkg/core/task"
	"bytetrade.io/web3os/installer/pkg/utils"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
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

// ~ CreateKsCore
type CreateKsCore struct {
	common.KubeAction
}

func (t *CreateKsCore) Execute(runtime connector.Runtime) error {
	masterNumIf, ok := t.PipelineCache.Get(common.CacheMasterNum)
	if !ok || masterNumIf == nil {
		return fmt.Errorf("failed to get master num")
	}

	kubeVersionIf, ok := t.PipelineCache.Get(common.CacheKubeletVersion)
	if !ok || kubeVersionIf == nil {
		return fmt.Errorf("failed to get kubelet version")
	}

	masterNum := masterNumIf.(int64)

	config, err := ctrl.GetConfig()
	if err != nil {
		return err
	}

	var appName = common.ChartNameKsCoreConfig
	var appPath = path.Join(runtime.GetFilesDir(), "apps", appName)

	actionConfig, settings, err := utils.InitConfig(config, common.NamespaceKubesphereSystem)
	if err != nil {
		return err
	}

	var values = make(map[string]interface{})
	values["Release"] = map[string]string{
		"Namespace": common.NamespaceKubesphereSystem,
	}

	if err := utils.InstallCharts(context.Background(), actionConfig, settings, appName,
		appPath, "", common.NamespaceKubesphereSystem, values); err != nil {
		logger.Errorf("failed to install %s chart: %v", appName, err)
		return err
	}

	appName = common.ChartNameKsCore
	appPath = path.Join(runtime.GetFilesDir(), "apps", appName)
	values = make(map[string]interface{})
	values["Release"] = map[string]string{
		"Namespace":    common.NamespaceKubesphereSystem,
		"ReplicaCount": fmt.Sprintf("%s", masterNum),
	}

	if err := utils.InstallCharts(context.Background(), actionConfig, settings, appName,
		appPath, "", common.NamespaceKubesphereSystem, values); err != nil {
		logger.Errorf("failed to install %s chart: %v", appName, err)
		return err
	}

	// manifestStr, err := util.Render(templates.KsCoreTempl, util.Data{
	// 	"MasterNum":   masterNum,
	// 	"KubeVersion": kubeVersionIf.(string),
	// })
	// if err != nil {
	// 	return err
	// }

	// var ksCoreConfigFileName = path.Join(runtime.GetFilesDir(), "apps", common.ChartNameKsCoreConfig, "values.yaml")
	// if err := ioutil.WriteFile(ksCoreConfigFileName, []byte(manifestStr), 0644); err != nil {
	// 	return errors.Wrapf(err, "failed to write ks-core-config values file %s", ksCoreConfigFileName)
	// }

	// var ksCoreFileName = path.Join(runtime.GetFilesDir(), "apps", common.ChartNameKsCore, "values.yaml")
	// if err := ioutil.WriteFile(ksCoreFileName, []byte(manifestStr), 0644); err != nil {
	// 	return errors.Wrapf(err, "failed to write ks-core values file %s", ksCoreFileName)
	// }

	return nil
}

// ~ CreateKsCoreManifests
type CreateKsCoreConfigManifests struct {
	common.KubeAction
}

func (t *CreateKsCoreConfigManifests) Execute(runtime connector.Runtime) error {
	var kscoreConfigCrdsPath = path.Join(runtime.GetFilesDir(), "apps", common.ChartNameKsCoreConfig, "crds")
	filepath.Walk(kscoreConfigCrdsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			_, err := runtime.GetRunner().SudoCmd(fmt.Sprintf("/usr/local/bin/kubectl apply -f %s", path), false, true)
			if err != nil {
				logger.Errorf("failed to apply %s: %v", path, err)
				return err
			}
		}
		return nil
	})

	return nil
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
		},
		Action:   new(PacthKsCore),
		Parallel: false,
		Retry:    0,
	}

	createKsCoreConfigManifests := &task.RemoteTask{
		Name:  "CreateKsCoreConfigManifests",
		Hosts: m.Runtime.GetHostsByRole(common.ETCD),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
			new(NotEqualDesiredVersion),
		},
		Action:   new(CreateKsCoreConfigManifests),
		Parallel: false,
		Retry:    0,
	}

	createKsCore := &task.RemoteTask{
		Name:  "CreateKsCore",
		Hosts: m.Runtime.GetHostsByRole(common.ETCD),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
			new(NotEqualDesiredVersion),
		},
		Action:   new(CreateKsCore),
		Parallel: false,
		Retry:    0,
	}

	m.Tasks = []task.Interface{
		checkKsCoreExist,
		pacthKsCore,
		createKsCoreConfigManifests,
		createKsCore,
	}
}
