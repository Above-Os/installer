package packages

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

// + 测试 full 包下载
func DownloadInstallPackage(kubeConf *common.KubeConf, path, version, arch string, pipelineCache *cache.Cache) error {
	installPackage := files.NewKubeBinary("full-package", arch, version, path, kubeConf.Arg.DownloadCommand)

	downloadFiles := []*files.KubeBinary{installPackage}
	filesMap := make(map[string]*files.KubeBinary)
	for _, downloadFile := range downloadFiles {
		if err := downloadFile.CreateBaseDir(); err != nil {
			return errors.Wrapf(errors.WithStack(err), "create file %s base dir failed", downloadFile.FileName)
		}

		filesMap[downloadFile.ID] = downloadFile
		if util.IsExist(downloadFile.Path()) {
			if err := downloadFile.SHA256Check(); err != nil {
				p := downloadFile.Path()
				_ = exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f %s", p)).Run()
			} else {
				logger.Infof("%s %s is existed", common.LocalHost, downloadFile.FileName)
				continue
			}
		}

		// todo doanload
		if err := downloadFile.Download(); err != nil {
			return fmt.Errorf("Failed to download %s binary: %s error: %w ", downloadFile.ID, downloadFile.GetCmd(), err)
		}
	}

	return nil
}

func DownloadPackage(kubeConf *common.KubeConf, path, version, arch string, pipelineCache *cache.Cache) error {
	// todo 这里会涉及多个文件的下载，且还会涉及 md5 的校验；同时还包括本地文件检查
	// file1 := files.NewKubeBinary("file1", arch, version, path, kubeConf.Arg.DownloadCommand)
	// file2 := files.NewKubeBinary("file2", arch, version, path, kubeConf.Arg.DownloadCommand)
	// file3 := files.NewKubeBinary("file3", arch, version, path, kubeConf.Arg.DownloadCommand)
	file4 := files.NewKubeBinary("kubekey", arch, "0.1.20", path, kubeConf.Arg.DownloadCommand) // todo test kubekey

	downloadFiles := []*files.KubeBinary{file4}

	filesMap := make(map[string]*files.KubeBinary)
	for _, downloadFile := range downloadFiles {
		if err := downloadFile.CreateBaseDir(); err != nil {
			return errors.Wrapf(errors.WithStack(err), "create file %s base dir failed", downloadFile.FileName)
		}

		// logger.Infof(common.LocalHost, "downloading %s %s %s ...", arch, downloadFile.ID, downloadFile.Version)

		filesMap[downloadFile.ID] = downloadFile
		if util.IsExist(downloadFile.Path()) {
			// download it again if it's incorrect
			if err := downloadFile.SHA256Check(); err != nil {
				p := downloadFile.Path()
				_ = exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -f %s", p)).Run()
			} else {
				logger.Infof(common.LocalHost, "%s is existed", downloadFile.ID)
				continue
			}
		}

		// todo
		if err := downloadFile.Download(); err != nil {
			return fmt.Errorf("Failed to download %s binary: %s error: %w ", downloadFile.ID, downloadFile.GetCmd(), err)
		}
	}

	return nil
}
