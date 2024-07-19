package storage

import (
	"fmt"
	"os/exec"

	kubekeyapiv1alpha2 "bytetrade.io/web3os/installer/apis/kubekey/v1alpha2"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/core/cache"
	cc "bytetrade.io/web3os/installer/pkg/core/common"
	corecommon "bytetrade.io/web3os/installer/pkg/core/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/task"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"bytetrade.io/web3os/installer/pkg/files"
	juicefsTemplates "bytetrade.io/web3os/installer/pkg/storage/templates"
	"bytetrade.io/web3os/installer/pkg/utils"
	"github.com/pkg/errors"
)

// - InstallJuiceFsModule
type InstallJuiceFsModule struct {
	common.KubeModule
}

func (m *InstallJuiceFsModule) Init() {
	m.Name = "InstallJuiceFs"

	downloadJuiceFs := &task.RemoteTask{
		Name:     "DownloadJuiceFs",
		Hosts:    m.Runtime.GetAllHosts(),
		Prepare:  &CheckJuiceFsExists{},
		Action:   new(DownloadJuiceFs),
		Parallel: false,
		Retry:    0,
	}

	installJuiceFs := &task.RemoteTask{
		Name:     "InstallJuiceFs",
		Hosts:    m.Runtime.GetAllHosts(),
		Prepare:  &CheckJuiceFsExists{},
		Action:   new(InstallJuiceFs),
		Parallel: false,
		Retry:    0,
	}

	enableJuiceFsService := &task.RemoteTask{
		Name:     "EnableJuiceFsService",
		Hosts:    m.Runtime.GetAllHosts(),
		Prepare:  &CheckJuiceFsExists{},
		Action:   new(EnableJuiceFsService),
		Parallel: false,
		Retry:    0,
	}

	m.Tasks = []task.Interface{
		downloadJuiceFs,
		installJuiceFs,
		enableJuiceFsService,
	}
}

// ~ CheckJuiceFsExists
type CheckJuiceFsExists struct {
	common.KubePrepare
}

func (p *CheckJuiceFsExists) PreCheck(runtime connector.Runtime) (bool, error) {
	if utils.IsExist(JuiceFsFile) {
		return false, nil
	}

	return true, nil
}

// ~ DownloadJuiceFs
type DownloadJuiceFs struct {
	common.KubeAction
}

func (t *DownloadJuiceFs) Execute(runtime connector.Runtime) error {
	binary := files.NewKubeBinary("juicefs", constants.OsArch, kubekeyapiv1alpha2.DefaultJuiceFsVersion, runtime.GetWorkDir())

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
		logger.Infof("%s downloading %s %s %s ...", common.LocalHost, runtime.RemoteHost().GetArch(), binary.ID, binary.Version)
		if err := binary.Download(); err != nil {
			return fmt.Errorf("Failed to download %s binary: %s error: %w ", binary.ID, binary.Url, err)
		}
	}

	t.ModuleCache.Set(common.CacheJuiceFsPath, binary.BaseDir)
	t.ModuleCache.Set(common.CacheJuiceFsFileName, binary.FileName)

	return nil
}

// ~ InstallJuiceFs
type InstallJuiceFs struct {
	common.KubeAction
}

func (t *InstallJuiceFs) Execute(runtime connector.Runtime) error {
	var redisPassword, _ = t.PipelineCache.GetMustString(common.CacheHostRedisPassword)
	var juiceFsBaseDir, _ = t.ModuleCache.GetMustString(common.CacheJuiceFsPath)
	var juiceFsFileName, _ = t.ModuleCache.GetMustString(common.CacheJuiceFsFileName)

	fmt.Println("---1---", juiceFsBaseDir)
	fmt.Println("---2---", juiceFsFileName)

	return fmt.Errorf("--------")

	if redisPassword == "" {
		return fmt.Errorf("redis password not found")
	}

	var cmd = fmt.Sprintf("cd %s && tar -zxf ./%s && chmod +x juicefs && install juicefs /usr/local/bin && install juicefs /sbin/mount.juicefs", juiceFsBaseDir, juiceFsFileName)
	if _, err := runtime.GetRunner().SudoCmdExt(cmd, false, true); err != nil {
		return err
	}

	var storageStr = getStorageTypeStr(t.PipelineCache, t.ModuleCache, t.KubeConf.Arg.Storage.StorageType)

	var redisService = fmt.Sprintf("redis://:%s@%s:6379/1", redisPassword, constants.LocalIp)
	cmd = fmt.Sprintf("%s format %s --storage %s", JuiceFsFile, redisService, t.KubeConf.Arg.Storage.StorageType)
	cmd = cmd + storageStr

	if _, err := runtime.GetRunner().SudoCmd(cmd, false, true); err != nil {
		return err
	}

	return nil
}

func getStorageTypeStr(pc, mc *cache.Cache, storageType string) string {
	var storageClusterId, _ = pc.GetMustString(common.CacheSTSClusterId)
	var formatStr string

	switch storageType {
	case common.Minio:
		formatStr = getMinioStr(pc)
	case common.OSS, common.S3:
		formatStr = getCloudStr(pc)
	}

	formatStr = formatStr + fmt.Sprintf(" %s --trash-days 0", storageClusterId)

	return formatStr
}

func getCloudStr(pc *cache.Cache) string {
	var storageBucket, _ = pc.GetMustString(common.CacheStorageBucket)
	var storageVendor, _ = pc.GetMustString(common.CacheStorageVendor)
	var storageAccessKey, _ = pc.GetMustString(common.CacheSTSAccessKey)
	var storageSecretKey, _ = pc.GetMustString(common.CacheSTSSecretKey)
	var storageToken, _ = pc.GetMustString(common.CacheSTSToken)

	var str = fmt.Sprintf(" --bucket %s", storageBucket)
	if storageVendor == "true" {
		if storageToken != "" {
			str = str + fmt.Sprintf(" --session-token %s", storageToken)
		}
	}
	if storageAccessKey != "" && storageSecretKey != "" {
		str = str + fmt.Sprintf(" --access-key %s --secret-key %s", storageAccessKey, storageSecretKey)
	}

	return str
}

func getMinioStr(pc *cache.Cache) string {
	var minioPassword, _ = pc.GetMustString(common.CacheMinioPassword)
	return fmt.Sprintf(" --bucket http://%s:9000/%s --access-key %s --secret-key -%s",
		constants.LocalIp, cc.TerminusDir, MinioRootUser, minioPassword)
}

// ~ EnableJuiceFsService
type EnableJuiceFsService struct {
	common.KubeAction
}

func (t *EnableJuiceFsService) Execute(runtime connector.Runtime) error {
	var redisPassword, _ = t.PipelineCache.GetMustString(common.CacheHostRedisPassword)
	var redisService = fmt.Sprintf("redis://:%s@%s:6379/1", redisPassword, constants.LocalIp)

	var data = util.Data{
		"JuiceFsBinPath":    JuiceFsFile,
		"JuiceFsCachePath":  JuiceFsCacheDir,
		"JuiceFsMetaDb":     redisService,
		"JuiceFsMountPoint": JuiceFsMountPointDir,
	}

	juiceFsServiceStr, err := util.Render(juicefsTemplates.JuicefsService, data)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), "render juicefs service template failed")
	}
	if err := util.WriteFile(JuiceFsServiceFile, []byte(juiceFsServiceStr), corecommon.FileMode0644); err != nil {
		return errors.Wrap(errors.WithStack(err), fmt.Sprintf("write juicefs service %s failed", JuiceFsServiceFile))
	}

	if _, err := runtime.GetRunner().SudoCmd("systemctl daemon-reload", false, false); err != nil {
		return err
	}

	if _, err := runtime.GetRunner().SudoCmd("systemctl restart juicefs", false, false); err != nil {
		return err
	}

	if _, err := runtime.GetRunner().SudoCmd("systemctl enable juicefs", false, false); err != nil {
		return err
	}

	if _, err := runtime.GetRunner().SudoCmd("systemctl --no-pager status juicefs", false, false); err != nil {
		return err
	}

	if _, err := runtime.GetRunner().SudoCmd(fmt.Sprintf("sleep 3 && test -d %s/.trash", JuiceFsMountPointDir), false, false); err != nil {
		return err
	}

	return nil
}
