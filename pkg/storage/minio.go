package storage

import (
	"fmt"
	"os/exec"
	"path"

	kubekeyapiv1alpha2 "bytetrade.io/web3os/installer/apis/kubekey/v1alpha2"
	"github.com/pkg/errors"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/constants"
	corecommon "bytetrade.io/web3os/installer/pkg/core/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/task"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"bytetrade.io/web3os/installer/pkg/files"
	minioTemplates "bytetrade.io/web3os/installer/pkg/storage/templates"
	"bytetrade.io/web3os/installer/pkg/utils"
)

// ~ CheckMinioState
// + todo debug
type CheckMinioState struct {
	common.KubeAction
}

func (t *CheckMinioState) Execute(runtime connector.Runtime) error {
	var cmd = "systemctl --no-pager status minio"
	stdout, err := runtime.GetRunner().SudoCmdExt(cmd, false, false)
	fmt.Println("---minio / 1---", stdout)
	fmt.Println("---minio / 2---", err)
	if err != nil {
		return fmt.Errorf("Minio Pending")
	}
	return nil
}

// ~ EnableMinio
type EnableMinio struct {
	common.KubeAction
}

func (t *EnableMinio) Execute(runtime connector.Runtime) error {
	var minioDataPath, _ = t.PipelineCache.GetMustString(common.CacheMinioDataPath)
	if minioDataPath == "" {
		return errors.New("get minio data path failed")
	}
	_, _ = runtime.GetRunner().SudoCmdExt("groupadd -r minio", false, false)
	_, _ = runtime.GetRunner().SudoCmdExt("useradd -M -r -g minio minio", false, false)
	_, _ = runtime.GetRunner().SudoCmdExt(fmt.Sprintf("chown minio:minio %s", minioDataPath), false, false)

	_, _ = runtime.GetRunner().SudoCmdExt("systemctl daemon-reload", false, false)
	_, _ = runtime.GetRunner().SudoCmd("systemctl restart minio", false, true)
	_, _ = runtime.GetRunner().SudoCmd("systemctl enable minio", false, true)
	_, _ = runtime.GetRunner().SudoCmd("systemctl --no-pager status minio", false, true)

	return nil
}

// ~ InstallMinio
type InstallMinio struct {
	common.KubeAction
}

func (t *InstallMinio) Execute(runtime connector.Runtime) error {
	minioPath, _ := t.PipelineCache.GetMustString(common.CacheMinioPath)
	if minioPath == "" {
		return errors.New("get minio path failed")
	}
	var minioDataPath = path.Join("terminus", "data", "minio", "vol1")

	if !utils.IsExist(minioDataPath) {
		utils.Mkdir(minioDataPath)
	}

	if _, err := runtime.GetRunner().SudoCmdExt(fmt.Sprintf("chmod +x %s", minioPath), false, false); err != nil {
		logger.Errorf("chmod +x %s error: %v", minioPath, err)
		return err
	}

	if _, err := runtime.GetRunner().SudoCmdExt(fmt.Sprintf("install %s /usr/local/bin", minioPath), false, false); err != nil {
		logger.Errorf("install minio %s error: %v", minioPath, err)
		return err
	}

	// write file
	var minioCmd = "/usr/local/bin/minio"
	var minioUser = "minioadmin"
	var minioPassword, _ = utils.GeneratePassword(16)
	var minioServiceFile = "/etc/systemd/system/minio.service"
	var minioEnvFile = "/etc/default/minio"
	var data = util.Data{
		"MinioCommand": minioCmd,
	}
	minioServiceStr, err := util.Render(minioTemplates.MinioService, data)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), "render minio service template failed")
	}
	if err := util.WriteFile(minioServiceFile, []byte(minioServiceStr), corecommon.FileMode0644); err != nil {
		return errors.Wrap(errors.WithStack(err), fmt.Sprintf("write minio service %s failed", minioServiceFile))
	}

	data = util.Data{
		"MinioDataPath": minioDataPath,
		"LocalIP":       constants.LocalIp,
		"User":          minioUser,
		"Password":      minioPassword,
	}
	minioEnvStr, err := util.Render(minioTemplates.MinioEnv, data)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), "render minio env template failed")
	}

	if err := util.WriteFile(minioEnvFile, []byte(minioEnvStr), corecommon.FileMode0644); err != nil {
		return errors.Wrap(errors.WithStack(err), fmt.Sprintf("write minio env %s failed", minioEnvFile))
	}

	t.PipelineCache.Set(common.CacheMinioDataPath, minioDataPath)
	t.PipelineCache.Set(common.CacheMinioPassword, minioPassword)

	return nil
}

// ~ DownloadMinio
type DownloadMinio struct {
	common.KubeAction
}

func (t *DownloadMinio) Execute(runtime connector.Runtime) error {
	var arch = constants.OsArch
	binary := files.NewKubeBinary("minio", arch, kubekeyapiv1alpha2.DefalutMultusVersion, runtime.GetWorkDir())

	if err := binary.CreateBaseDir(); err != nil {
		return errors.Wrapf(errors.WithStack(err), "create file %s base dir failed", binary.FileName)
	}

	var exists = util.IsExist(binary.Path())
	if exists {
		p := binary.Path()
		if err := binary.SHA256Check(); err != nil {
			_ = exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f %s", p)).Run()
		} else {
			return nil
		}
	}

	if !exists || binary.OverWrite {
		logger.Infof("%s downloading %s %s %s ...", common.LocalHost, arch, binary.ID, binary.Version)
		if err := binary.Download(); err != nil {
			return fmt.Errorf("Failed to download %s binary: %s error: %w ", binary.ID, binary.Url, err)
		}
	}

	t.PipelineCache.Set(common.CacheMinioPath, binary.Path())

	return nil
}

// - InstallMinioModule
type InstallMinioModule struct {
	common.KubeModule
}

func (m *InstallMinioModule) Init() {
	m.Name = "InstallMinio"

	downloadMinio := &task.RemoteTask{
		Name:     "Download",
		Hosts:    m.Runtime.GetAllHosts(),
		Action:   &DownloadMinio{},
		Parallel: false,
	}

	installMinio := &task.RemoteTask{
		Name:     "Install",
		Hosts:    m.Runtime.GetAllHosts(),
		Action:   &InstallMinio{},
		Parallel: false,
	}

	enableMinio := &task.RemoteTask{
		Name:     "Enable",
		Hosts:    m.Runtime.GetAllHosts(),
		Action:   &EnableMinio{},
		Parallel: false,
	}

	checkMinioState := &task.RemoteTask{
		Name:     "CheckMinioState",
		Hosts:    m.Runtime.GetAllHosts(),
		Action:   &CheckMinioState{},
		Parallel: false,
		Retry:    60,
		Delay:    5,
	}

	m.Tasks = []task.Interface{
		downloadMinio,
		installMinio,
		enableMinio,
		checkMinioState,
	}
}

// - InstallMinioClusterModule
type InstallMinioClusterModule struct {
	common.KubeModule
}

func (m *InstallMinioClusterModule) Init() {
	m.Name = "InstallMinioCluster"
}

// ~ InstallMinioOperator
type InstallMinioOperator struct {
	common.KubeAction
}

func (t *InstallMinioOperator) Execute(runtime connector.Runtime) error {
	var arch = constants.OsArch
	binary := files.NewKubeBinary("minio-operator", arch, kubekeyapiv1alpha2.DefaultMinioOperatorVersion, runtime.GetWorkDir())

	if err := binary.CreateBaseDir(); err != nil {
		return errors.Wrapf(errors.WithStack(err), "create file %s base dir failed", binary.FileName)
	}

	var exists = util.IsExist(binary.Path())
	if exists {
		p := binary.Path()
		if err := binary.SHA256Check(); err != nil {
			_ = exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f %s", p)).Run()
		} else {
			return nil
		}
	}

	if !exists || binary.OverWrite {
		logger.Infof("%s downloading %s %s %s ...", common.LocalHost, arch, binary.ID, binary.Version)
		if err := binary.Download(); err != nil {
			return fmt.Errorf("Failed to download %s binary: %s error: %w ", binary.ID, binary.Url, err)
		}
	}

	_, _ = runtime.GetRunner().SudoCmd(fmt.Sprintf("tar zxvf %s", binary.Path()), false, true)
	_, _ = runtime.GetRunner().SudoCmd(fmt.Sprintf("install -m 755 %s/minio-operator /usr/local/bin/minio-operator", binary.BaseDir), false, true)

	var minioData, _ = t.PipelineCache.GetMustString(common.CacheMinioDataPath)
	var minioPassword, _ = t.PipelineCache.GetMustString(common.CacheMinioPassword)
	var cmd = fmt.Sprintf("/usr/local/bin/minio-operator init --address %s --cafile /etc/ssl/etcd/ssl/ca.pem --certfile /etc/ssl/etcd/ssl/node-%s.pem --keyfile /etc/ssl/etcd/ssl/node-%s-key.pem --volume %s --password %s",
		constants.LocalIp, runtime.RemoteHost().GetName(), runtime.RemoteHost().GetName(), minioData, minioPassword)

	_, _ = runtime.GetRunner().SudoCmd(cmd, false, true)

	return nil
}
