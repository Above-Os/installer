package storage

import (
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/connector"
)

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
