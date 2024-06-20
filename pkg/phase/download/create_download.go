package download

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/module"
	"bytetrade.io/web3os/installer/pkg/core/pipeline"
	"bytetrade.io/web3os/installer/pkg/packages"
)

func NewPackageDownloadPipeline(runtime *common.KubeRuntime) error {
	m := []module.Module{
		&packages.PackagesModule{},
	}

	p := pipeline.Pipeline{
		Name:    "CreatePackageDownloadPipeline",
		Modules: m,
		Runtime: runtime,
	}

	go p.Start()

	return nil
}

// todo 这里应该是带参数的；但也可能是从数据库中查询配置信息，然后来下载指定的 package
func CreateDownload(args common.Argument, downloadCmd string) error {
	args.DownloadCommand = func(path, url string) string {
		// this is an extension point for downloading tools, for example users can set the timeout, proxy or retry under
		// some poor network environment. Or users even can choose another cli, it might be wget.
		// perhaps we should have a build-in download function instead of totally rely on the external one
		return fmt.Sprintf(downloadCmd, path, url)
	}

	runtime, err := common.NewKubeRuntime("", args)
	if err != nil {
		return err
	}

	if err := NewPackageDownloadPipeline(runtime); err != nil {
		return err
	}

	return nil
}
