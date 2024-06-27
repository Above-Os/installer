package model

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type InstallModelReq struct {
	DomainName string `json:"terminus_os_domainname"`
	UserName   string `json:"terminus_os_username" validate:"required"`
	KubeType   string `json:"kube_type" validate:"kubeTypeValid"`
	Vendor     string `json:"vendor"`
	GpuEnable  int    `json:"gpu_enable" validate:"oneof=0 1"`
	GpuShare   int    `json:"gpu_share" validate:"required_with=GpuEnable,oneof=0 1"`
	Version    string `json:"version"`
}

func KubeTypeValid(fl validator.FieldLevel) bool {
	kubeType := strings.ToLower(fl.Field().String())
	return kubeType == "" || kubeType == "k3s" || kubeType == "k8s"
}
