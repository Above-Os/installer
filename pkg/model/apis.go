package model

type InstallModelReq struct {
	DomainName     string `json:"terminus_os_domainname,"`
	UserName       string `json:"terminus_os_username"`
	KubeType       string `json:"kube_type"`
	Vendor         string `json:"vendor"`
	GpuEnable      int    `json:"gpu_enable"`
	GpuShare       int    `json:"gpu_share"`
	Version        string `json:"version"`
	DownloadImages int    `json:"download_images"`
	DownloadDeps   int    `json:"download_deps"`
}
