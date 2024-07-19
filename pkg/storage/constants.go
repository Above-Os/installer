package storage

import (
	"path"

	cc "bytetrade.io/web3os/installer/pkg/core/common"
)

var (
	RedisRootDir             = path.Join(cc.TerminusDir, "data", "redis")
	RedisConfigDir           = path.Join(RedisRootDir, "etc")
	RedisDataDir             = path.Join(RedisRootDir, "data")
	RedisLogDir              = path.Join(RedisRootDir, "log")
	RedisRunDir              = path.Join(RedisRootDir, "run")
	RedisConfigFile          = path.Join(RedisConfigDir, "redis.conf")
	RedisServiceFile         = path.Join("etc", "systemd", "system", "redis-server.service")
	RedisServerFile          = path.Join("usr", "bin", "redis-server")
	RedisCliFile             = path.Join("usr", "bin", "redis-cli")
	RedisServerInstalledFile = path.Join("usr", "local", "bin", "redis-server")
	RedisCliInstalledFile    = path.Join("usr", "local", "bin", "redis-cli")

	JuiceFsFile          = path.Join("usr", "local", "bin", "juicefs")
	JuiceFsCacheDir      = path.Join(cc.TerminusDir, "jfscache")
	JuiceFsMountPointDir = path.Join(cc.TerminusDir, "rootfs")
	JuiceFsServiceFile   = path.Join("etc", "systemd", "system", "juicefs.service")

	MinioRootUser    = "minioadmin"
	MinioDataDir     = path.Join(cc.TerminusDir, "data", "minio", "vol1")
	MinioFile        = path.Join("usr", "local", "bin", "minio")
	MinioServiceFile = path.Join("etc", "systemd", "system", "minio.service")
	MinioConfigFile  = path.Join("etc", "default", "minio")

	MinioOperatorFile = path.Join("usr", "local", "bin", "minio-operator")
)
