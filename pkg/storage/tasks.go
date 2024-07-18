package storage

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/model"
)

// ~ StopJuiceFS
type StopJuiceFS struct {
	common.KubeAction
}

func (t *StopJuiceFS) Execute(runtime connector.Runtime) error {
	var cmd = "systemctl stop juicefs; systemctl disable juicefs"
	stdout, err := runtime.GetRunner().SudoCmdExt(cmd, false, false)
	fmt.Println("---juiceFS / 1---", stdout)
	fmt.Println("---juiceFS / 2---", err)

	cmd = "rm -rf /var/jfsCache /terminus/jfscache"
	stdout, err = runtime.GetRunner().SudoCmdExt(cmd, false, false)
	fmt.Println("---juiceFS / 3---", stdout)
	fmt.Println("---juiceFS / 4---", err)

	return nil
}

// ~ StopMinio
type StopMinio struct {
	common.KubeAction
}

func (t *StopMinio) Execute(runtime connector.Runtime) error {
	var cmd = "systemctl stop minio; systemctl disable minio"
	stdout, err := runtime.GetRunner().SudoCmdExt(cmd, false, false)
	fmt.Println("---minio / 1---", stdout)
	fmt.Println("---minio / 2---", err)
	return nil
}

// ~ StopMinioOperator
type StopMinioOperator struct {
	common.KubeAction
}

func (t *StopMinioOperator) Execute(runtime connector.Runtime) error {
	var cmd = "systemctl stop minio-operator; systemctl disable minio-operator"
	stdout, err := runtime.GetRunner().SudoCmdExt(cmd, false, false)
	fmt.Println("---minio-operator / 1---", stdout)
	fmt.Println("---minio-operator / 2---", err)
	return nil
}

// ~ StopRedis
type StopRedis struct {
	common.KubeAction
}

func (t *StopRedis) Execute(runtime connector.Runtime) error {
	var cmd = "systemctl stop redis-server; systemctl disable redis-server"
	stdout, err := runtime.GetRunner().SudoCmdExt(cmd, false, false)
	fmt.Println("---redis / 1---", stdout)
	fmt.Println("---redis / 2---", err)

	stdout, err = runtime.GetRunner().SudoCmd("killall -9 redis-server", false, true)
	fmt.Println("---redis / 3---", stdout)
	fmt.Println("---redis / 4---", err)

	stdout, err = runtime.GetRunner().SudoCmd("unlink /usr/bin/redis-server; unlink /usr/bin/redis-cli", false, true)
	fmt.Println("---redis / 5---", stdout)
	fmt.Println("---redis / 6---", err)

	return nil
}

// ~ RemoveTerminusFiles
type RemoveTerminusFiles struct {
	common.KubeAction
}

func (t *RemoveTerminusFiles) Execute(runtime connector.Runtime) error {
	var files = []string{
		"/usr/local/bin/redis-*",
		"/usr/bin/redis-*",
		"/sbin/mount.juicefs",
		"/etc/init.d/redis-server",
		"/usr/local/bin/juicefs",
		"/usr/local/bin/minio",
		"/usr/local/bin/velero",
		"/etc/systemd/system/redis-server.service",
		"/etc/systemd/system/minio.service",
		"/etc/systemd/system/minio-operator.service",
		"/etc/systemd/system/juicefs.service",
		"/etc/systemd/system/containerd.service",
		"/terminus/",
	}

	for _, f := range files {
		runtime.GetRunner().SudoCmdExt(fmt.Sprintf("rm -rf %s", f), false, true)
	}

	return nil
}

// +

// ~ task SaveInstallConfigTask
type SaveInstallConfigTask struct {
	common.KubeAction
}

func (t *SaveInstallConfigTask) Execute(runtime connector.Runtime) error {
	var installReq model.InstallModelReq
	var ok bool
	if installReq, ok = any(t.KubeConf.Arg.Request).(model.InstallModelReq); !ok {
		return fmt.Errorf("invalid install model req %+v", t.KubeConf.Arg.Request)
	}

	return t.KubeConf.Arg.Provider.SaveInstallConfig(installReq)
}
