package download

import (
	"bytetrade.io/web3os/installer/pkg/common"
)

func DownloadMiniInstallPackage(args common.Argument) error {
	return nil
}

// todo 这里应该是带参数的；但也可能是从数据库中查询配置信息，然后来下载指定的 package
func CreateDownload(args common.Argument) error {
	// runtime, err := common.NewKubeRuntime("", args)
	// if err != nil {
	// 	return err
	// }

	// m := []module.Module{
	// 	&packages.PackagesModule{},
	// }

	// p := pipeline.Pipeline{
	// 	Name:    "CreatePackageDownloadPipeline",
	// 	Modules: m,
	// 	Runtime: runtime,
	// }

	// go p.Start()

	return nil
}
