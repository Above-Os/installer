package v1

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/api/response"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/storage"
	"bytetrade.io/web3os/installer/pkg/phase/download"
	"bytetrade.io/web3os/installer/pkg/phase/mock"
	"bytetrade.io/web3os/installer/pkg/pipelines"
	"github.com/emicklei/go-restful/v3"
)

type Handler struct {
	// apis.Base
	// appService *app_service.Client
	StorageProvider storage.Provider
}

func New() *Handler {
	// as := app_service.NewAppServiceClient()
	// return &Handler{
	// 	appService: as,
	// }
	return &Handler{}
}

func (h *Handler) handlerTest(req *restful.Request, resp *restful.Response) {
	logger.Infof("handler test req: %s", req.Request.Method)
	response.SuccessNoData(resp)
}

// + test func
func (h *Handler) handlerGreetings(req *restful.Request, resp *restful.Response) {
	logger.Infof("handler greetings req: %s", req.Request.Method)

	if err := mock.Greetings(); err != nil {
		logger.Errorf("greetings failed %v", err)
	}

	response.SuccessNoData(resp)
}

func (h *Handler) handlerDownload(req *restful.Request, resp *restful.Response) {
	logger.Infof("handler download req: %s", req.Request.Method)

	arg := common.Argument{}

	if err := download.CreateDownload(arg, "curl -L -o %s %s"); err != nil {
		logger.Errorf("download failed %v", err)
	}

	response.SuccessNoData(resp)
}

// ~ 测试安装 kk
func (h *Handler) handlerInstallKk(req *restful.Request, resp *restful.Response) {
	logger.Infof("handler installer req: %s", req.Request.Method)

	arg := common.Argument{}

	if err := pipelines.InstallKubekeyPipeline(arg); err != nil {
		fmt.Println("---api installer kk / err---", err)
	}

	response.SuccessNoData(resp)
}

// todo 一个完整的测试流程，下载 full 包并安装
func (h *Handler) handlerInstallTerminus(req *restful.Request, resp *restful.Response) {
	logger.Infof("handler installer req: %s", req.Request.Method)

	arg := common.Argument{}
	if err := pipelines.InstallTerminusPipeline(arg); err != nil {
		fmt.Println("---api installer terminus / err---", err)
	}

	response.SuccessNoData(resp)
}
