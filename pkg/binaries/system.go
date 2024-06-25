package binaries

import (
	"fmt"
	"os/exec"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/cache"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/util"
	"bytetrade.io/web3os/installer/pkg/files"
	"github.com/pkg/errors"
)

func DownloadUbutun24AppArmor(path, version, arch string, pipelineCache *cache.Cache) error {
	apparmor := files.NewKubeBinary("apparmor", arch, version, path, nil)

	if err := apparmor.CreateBaseDir(); err != nil {
		return errors.Wrapf(errors.WithStack(err), "create file %s base dir failed", apparmor.FileName)
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
		return fmt.Errorf("Failed to download %s binary: %s error: %w ", apparmor.ID, apparmor.GetCmd(), err)
	}

	binariesMap := make(map[string]*files.KubeBinary)
	binariesMap[apparmor.ID] = apparmor
	pipelineCache.Set(common.KubeBinaries+"-"+arch, binariesMap)
	return nil
}
