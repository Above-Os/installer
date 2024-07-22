package v1

import (
	"fmt"
	"os"

	"bytetrade.io/web3os/installer/pkg/api/response"
	"bytetrade.io/web3os/installer/pkg/common"
	corecommon "bytetrade.io/web3os/installer/pkg/core/common"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/model"
	"bytetrade.io/web3os/installer/pkg/phase/mock"
	"bytetrade.io/web3os/installer/pkg/pipelines"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	// apis.Base
	// appService *app_service.Client
	validate *validator.Validate
}

func New() *Handler {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("kubeTypeValid", model.KubeTypeValid)
	return &Handler{
		validate: v,
	}
}

// + install
// ~ 先下载完整包进行安装，需要提取日志写入数据库
func (h *Handler) handlerInstall(req *restful.Request, resp *restful.Response) {
	logger.Infof("handler install req: %s", req.Request.Method)

	var reqModel model.InstallModelReq
	err := req.ReadEntity(&reqModel)
	if err != nil {
		response.HandleError(resp, err)
		return
	}

	if err = h.validate.Struct(&reqModel); err != nil {
		if validationErrors := err.(validator.ValidationErrors); validationErrors != nil {
			logger.Errorf("handler install request parameter invalid: %v", validationErrors)
			response.HandleError(resp, fmt.Errorf("handler install request parameter invalid"))
			return
		}
	}

	if reqModel.Config.DomainName == "" {
		reqModel.Config.DomainName = corecommon.DefaultDomainName
	}

	arg := common.Argument{
		KsEnable:         true,
		KsVersion:        common.DefaultKubeSphereVersion,
		Request:          reqModel,
		InstallPackages:  false,
		SKipPushImages:   false,
		ContainerManager: common.Containerd,
		RegistryMirrors:  GetEnv("REGISTRY_MIRRORS", reqModel.Config.RegistryMirrors),
		Proxy:            GetEnv("PROXY", reqModel.Config.Proxy),
	}

	switch reqModel.Config.KubeType {
	case common.K3s:
		arg.KubernetesVersion = common.DefaultK3sVersion
	case common.K8s:
		arg.KubernetesVersion = common.DefaultK8sVersion
	}

	if err := pipelines.InstallTerminusPipeline(arg); err != nil {
		response.HandleError(resp, err)
		return
	}

	response.SuccessNoData(resp)
}

func (h *Handler) handlerStatus(req *restful.Request, resp *restful.Response) {
	var timespan = req.QueryParameter("time")
	if timespan == "" {
		timespan = "0"
	}

	response.Success(resp, "ok")
}

func (h *Handler) handlerGreetings(req *restful.Request, resp *restful.Response) {
	logger.Infof("handler greetings req: %s", req.Request.Method)

	if err := mock.Greetings(); err != nil {
		logger.Errorf("greetings failed %v", err)
	}

	response.SuccessNoData(resp)
}

// todo 一个完整的测试流程，下载 full 包并安装
func (h *Handler) handlerInstallTerminus(req *restful.Request, resp *restful.Response) {
	// logger.Infof("handler installer req: %s", req.Request.Method)

	// arg := common.Argument{}
	// if err := pipelines.InstallTerminusPipeline(arg); err != nil {
	// 	fmt.Println("---api installer terminus / err---", err)
	// }

	response.SuccessNoData(resp)
}

func GetEnv(key string, arg string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return arg
}
