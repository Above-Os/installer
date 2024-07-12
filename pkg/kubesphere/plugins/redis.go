package plugins

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/prepare"
	"bytetrade.io/web3os/installer/pkg/core/task"
	"github.com/pkg/errors"
	kubeErr "k8s.io/apimachinery/pkg/api/errors"
)

// ~ DeployRedis
type DeployRedis struct {
	common.KubeAction
}

func (t *DeployRedis) Execute(runtime connector.Runtime) error {
	redisPwd, ok := t.ModuleCache.Get(common.CacheRedisPassword)
	if !ok {
		return fmt.Errorf("get redis password from module cache failed")
	}

	if _, err := runtime.GetRunner().SudoCmd(fmt.Sprintf("/usr/local/bin/kubectl -n %s create secret generic redis-secret --from-literal=auth=%s", common.NamespaceKubesphereSystem, redisPwd), false, true); err != nil {
		if !kubeErr.IsAlreadyExists(err) {
			return errors.Wrap(errors.WithStack(err), "create redis secret failed")
		}
	}

	rver, _ := runtime.GetRunner().SudoCmd(fmt.Sprintf("/usr/local/bin/kubectl get pod -n %s -l app=%s,tier=database,version=%s-4.0 | wc -l",
		common.NamespaceKubesphereSystem, common.NamespaceKubesphereSystem), false, true)

	fmt.Println("---rver---", rver)
	if rver != "0" {
		var cmd = fmt.Sprintf("/usr/local/bin/kubectl get svc -n %s %s -o yaml > %s/redis-svc-backup.yaml && /usr/local/bin/kubectl delete svc -n %s %s",
			common.NamespaceKubesphereSystem, common.ChartNameRedis, common.KubeManifestDir, common.NamespaceKubesphereSystem, common.ChartNameRedis)
		if _, err := runtime.GetRunner().SudoCmd(cmd, false, true); err != nil {
			logger.Errorf("failed to backup %s svc: %v", common.ChartNameRedis, err)
		}
	}

	// todo  enableHA

	// todo deploying redis
	return nil
}

// ~ DeployRedisModule
type DeployRedisModule struct {
	common.KubeModule
}

func (m *DeployRedisModule) Init() {
	m.Name = "DeployRedis"

	createRedis := &task.RemoteTask{
		Name:  "DeployRedis",
		Hosts: m.Runtime.GetHostsByRole(common.Master),
		Prepare: &prepare.PrepareCollection{
			new(common.OnlyFirstMaster),
			new(NotEqualDesiredVersion),
			new(GenerateRedisPassword),
		},
		Action:   new(DeployRedis),
		Parallel: false,
	}

	m.Tasks = []task.Interface{
		createRedis,
	}
}
