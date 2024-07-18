package storage

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
	"bytetrade.io/web3os/installer/pkg/utils"
)

// ~ CheckEtcdSSL
type CheckEtcdSSL struct {
	common.KubePrepare
}

func (p *CheckEtcdSSL) PreCheck(runtime connector.Runtime) (bool, error) {
	var files = []string{
		"/etc/ssl/etcd/ssl/ca.pem",
		fmt.Sprintf("/etc/ssl/etcd/ssl/node-%s-key.pem", runtime.RemoteHost().GetName()),
		fmt.Sprintf("/etc/ssl/etcd/ssl/node-%s.pem", runtime.RemoteHost().GetName()),
	}
	for _, f := range files {
		if !utils.IsExist(f) {
			return false, nil
		}
	}
	return true, nil
}

// ~ CheckStorageType
type CheckStorageType struct {
	common.KubePrepare
	StorageType string
}

func (p *CheckStorageType) PreCheck(runtime connector.Runtime) (bool, error) {
	storageType, _ := p.PipelineCache.GetMustString(common.CacheStorageType)
	if storageType == "" || storageType != p.StorageType {
		return false, nil
	}
	return true, nil
}

// ~ CheckStorageVendor
type CheckStorageVendor struct {
	common.KubePrepare
}

func (p *CheckStorageVendor) PreCheck(runtime connector.Runtime) (bool, error) {
	storageVendor, _ := p.PipelineCache.GetMustString(common.CacheStorageVendor)
	if storageVendor != "true" {
		return false, nil
	}

	if storageType, _ := p.PipelineCache.GetMustString(common.CacheStorageType); storageType != "s3" && storageType != "oss" {
		return false, nil
	}

	if storageBucket, _ := p.PipelineCache.GetMustString(common.CacheStorageBucket); storageBucket == "" {
		return false, nil
	}

	return true, nil
}
