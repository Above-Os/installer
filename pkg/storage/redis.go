package storage

import (
	"fmt"
	"os/exec"
	"path"

	kubekeyapiv1alpha2 "bytetrade.io/web3os/installer/apis/kubekey/v1alpha2"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/constants"
	corecommon "bytetrade.io/web3os/installer/pkg/core/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"bytetrade.io/web3os/installer/pkg/files"
	"bytetrade.io/web3os/installer/pkg/utils"

	redisTemplates "bytetrade.io/web3os/installer/pkg/storage/templates"
	"github.com/pkg/errors"
)

// - InstallRedisModule
type InstallRedisModule struct {
	common.KubeModule
}

func (m *InstallRedisModule) Init() {
	m.Name = "InstallRedis"
}

// ~ InstallRedis
type InstallRedis struct {
	common.KubeAction
}

func (t *InstallRedis) Execute(runtime connector.Runtime) error {
	binary := files.NewKubeBinary("redis", constants.OsArch, kubekeyapiv1alpha2.DefaultRedisVersion, runtime.GetWorkDir())

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

	if _, err := runtime.GetRunner().SudoCmdExt(fmt.Sprintf("cp %s && tar xf ./%s", binary.BaseDir, binary.FileName), false, false); err != nil {
		return errors.Wrapf(errors.WithStack(err), "untar redis failed")
	}

	var cmd = fmt.Sprintf("cd %s/redis-%s && make -j%d && make install", binary.BaseDir, binary.Version, constants.CpuPhysicalCount)
	if _, err := runtime.GetRunner().SudoCmdExt(cmd, false, true); err != nil {
		return err
	}
	if _, err := runtime.GetRunner().SudoCmdExt("ln -s /usr/local/bin/redis-server /usr/bin/redis-server", false, true); err != nil {
		return err
	}
	if _, err := runtime.GetRunner().SudoCmdExt("ln -s /usr/local/bin/redis-cli /usr/bin/redis-cli", false, true); err != nil {
		return err
	}

	// mkdir
	var redisRootPath = path.Join("terminus", "data", "redis")
	var redisPassword, _ = utils.GeneratePassword(16)
	if !utils.IsExist(redisRootPath) {
		utils.Mkdir(fmt.Sprintf("%s/etc", redisRootPath))
		utils.Mkdir(fmt.Sprintf("%s/data", redisRootPath))
		utils.Mkdir(fmt.Sprintf("%s/log", redisRootPath))
		utils.Mkdir(fmt.Sprintf("%s/run", redisRootPath))
	}

	// config
	var redisConfFile = path.Join(redisRootPath, "etc", "redis.conf")
	var data = util.Data{
		"LocalIP":  constants.LocalIp,
		"RootPath": redisRootPath,
		"Password": redisPassword,
	}
	redisConfStr, err := util.Render(redisTemplates.RedisConf, data)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), "render redis conf template failed")
	}
	if err := util.WriteFile(redisConfFile, []byte(redisConfStr), corecommon.FileMode0600); err != nil {
		return errors.Wrap(errors.WithStack(err), fmt.Sprintf("write redis conf %s failed", redisConfFile))
	}

	var redisServiceFile = path.Join("etc", "systemd", "system", "redis-server.service")
	data = util.Data{
		"RedisBinPath":  "/usr/bin/redis-server",
		"RootPath":      redisRootPath,
		"RedisConfPath": redisConfFile,
	}
	redisServiceStr, err := util.Render(redisTemplates.RedisService, data)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), "render redis service template failed")
	}
	if err := util.WriteFile(redisServiceFile, []byte(redisServiceStr), corecommon.FileMode0644); err != nil {
		return errors.Wrap(errors.WithStack(err), fmt.Sprintf("write redis service %s failed", redisServiceFile))
	}

	if _, err := runtime.GetRunner().SudoCmdExt("sysctl -w vm.overcommit_memory=1 net.core.somaxconn=10240", false, false); err != nil {
		return err
	}
	if _, err := runtime.GetRunner().SudoCmdExt("systemctl daemon-reload", false, false); err != nil {
		return err
	}
	if _, err := runtime.GetRunner().SudoCmdExt("systemctl restart redis-server", false, false); err != nil {
		return err
	}
	if _, err := runtime.GetRunner().SudoCmdExt("systemctl enable redis-server", false, false); err != nil {
		return err
	}

	cmd = "( sleep 10 && systemctl --no-pager status redis-server ) || ( systemctl restart redis-server && sleep 3 && systemctl --no-pager status redis-server ) || ( systemctl restart redis-server && sleep 3 && systemctl --no-pager status redis-server )"
	if _, err = runtime.GetRunner().SudoCmdExt(cmd, false, false); err != nil {
		return err
	}

	cmd = fmt.Sprintf("awk '/requirepass/{print \\$NF}' %s", redisConfFile)
	rpwd, _ := runtime.GetRunner().SudoCmdExt(cmd, false, false)
	if rpwd == "" {
		return fmt.Errorf("get redis password failed")
	}

	cmd = fmt.Sprintf("/usr/bin/redis-cli -h %s -a %s ping", constants.LocalIp, rpwd)
	if pong, _ := runtime.GetRunner().SudoCmdExt(cmd, false, false); pong != "PONG" {
		return fmt.Errorf("failed to connect redis server: %s:6379", constants.LocalIp)
	}

	t.PipelineCache.Set(common.CacheHostRedisPassword, redisPassword)

	return nil
}
