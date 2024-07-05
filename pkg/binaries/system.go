package binaries

import (
	"fmt"
	"os/exec"
	"path"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/cache"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"bytetrade.io/web3os/installer/pkg/files"
	"github.com/pkg/errors"
)

func DownloadSocat(path, version, arch string, pipelineCache *cache.Cache) (string, string, error) {
	socat := files.NewKubeBinary("socat", arch, version, path)

	if err := socat.CreateBaseDir(); err != nil {
		return "", "", errors.Wrapf(errors.WithStack(err), "create file %s base dir failed", socat.FileName)
	}

	logger.Infof("%s downloading %s %s ...", common.LocalHost, socat.ID, socat.Version)

	if util.IsExist(socat.Path()) {
		// download it again if it's incorrect
		if err := socat.SHA256Check(); err != nil {
			p := socat.Path()
			_ = exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f %s", p)).Run()
		} else {
			logger.Infof("%s %s is existed", common.LocalHost, socat.ID)

		}
	}

	if err := socat.Download(); err != nil {
		return "", "", fmt.Errorf("Failed to download %s binary: %s error: %w ", socat.ID, socat.Url, err)
	}

	binariesMap := make(map[string]*files.KubeBinary)
	binariesMap[socat.ID] = socat
	pipelineCache.Set(common.KubeBinaries+"-"+arch, binariesMap)
	return socat.BaseDir, socat.FileName, nil
}

func DownloadConntrack(path, version, arch string, pipelineCache *cache.Cache) (string, string, error) {
	conntrack := files.NewKubeBinary("conntrack", arch, version, path)

	if err := conntrack.CreateBaseDir(); err != nil {
		return "", "", errors.Wrapf(errors.WithStack(err), "create file %s base dir failed", conntrack.FileName)
	}

	logger.Infof("%s downloading %s %s ...", common.LocalHost, conntrack.ID, conntrack.Version)

	if util.IsExist(conntrack.Path()) {
		// download it again if it's incorrect
		if err := conntrack.SHA256Check(); err != nil {
			p := conntrack.Path()
			_ = exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f %s", p)).Run()
		} else {
			logger.Infof("%s %s is existed", common.LocalHost, conntrack.ID)

		}
	}

	if err := conntrack.Download(); err != nil {
		return "", "", fmt.Errorf("Failed to download %s binary: %s error: %w ", conntrack.ID, conntrack.Url, err)
	}

	binariesMap := make(map[string]*files.KubeBinary)
	binariesMap[conntrack.ID] = conntrack
	pipelineCache.Set(common.KubeBinaries+"-"+arch, binariesMap)
	return conntrack.BaseDir, conntrack.FileName, nil
}

func DownloadUbutun24AppArmor(prePath, version, arch string, pipelineCache *cache.Cache) (string, error) {
	apparmor := files.NewKubeBinary("apparmor", arch, version, prePath)

	if err := apparmor.CreateBaseDir(); err != nil {
		return "", errors.Wrapf(errors.WithStack(err), "create file %s base dir failed", apparmor.FileName)
	}

	logger.Infof("%s downloading %s %s %s ...", common.LocalHost, arch, apparmor.ID, apparmor.Version)

	if util.IsExist(apparmor.Path()) {
		// download it again if it's incorrect
		if err := apparmor.SHA256Check(); err != nil {
			p := apparmor.Path()
			_ = exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f %s", p)).Run()
		} else {
			logger.Infof("%s %s is existed", common.LocalHost, apparmor.ID)

		}
	}

	if err := apparmor.Download(); err != nil {
		return "", fmt.Errorf("Failed to download %s binary: %s error: %w ", apparmor.ID, apparmor.Url, err)
	}

	binariesMap := make(map[string]*files.KubeBinary)
	binariesMap[apparmor.ID] = apparmor
	pipelineCache.Set(common.KubeBinaries+"-"+arch, binariesMap)
	return path.Join(apparmor.BaseDir, apparmor.FileName), nil

}
