package v1

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/api/response"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/core/logger"
	"bytetrade.io/web3os/installer/pkg/core/storage"
	"bytetrade.io/web3os/installer/pkg/model"
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

// ~ get public ip
func (h *Handler) handlerPublicIp(req *restful.Request, resp *restful.Response) {
	var data = make(map[string]interface{})
	data["public_ip"] = "13.92.32.12"

	if 1 == 1 {
		response.HandleError(resp, fmt.Errorf("app %s entrances not found", "TEST"))
		return
	}

	response.Success(resp, data)
}

func (h *Handler) handlerConfig(req *restful.Request, resp *restful.Response) {
	response.SuccessNoData(resp)
}

// ~ install
func (h *Handler) handlerInstall(req *restful.Request, resp *restful.Response) {
	var reqModel model.InstallModelReq
	err := req.ReadEntity(&reqModel)
	if err != nil {
		response.HandleError(resp, err)
		return
	}

	// todo 写入数据库

}

func (h *Handler) handlerProgress(req *restful.Request, resp *restful.Response) {

}

func (h *Handler) handlerStatus(req *restful.Request, resp *restful.Response) {

}

// - test func
func (h *Handler) handlerTest(req *restful.Request, resp *restful.Response) {
	logger.Infof("handler test req: %s", req.Request.Method)
	response.SuccessNoData(resp)
}

// + 测试安装的接口
func (h *Handler) handlerInst(req *restful.Request, resp *restful.Response) {
	args := common.Argument{}
	runtime, err := common.NewKubeRuntime(common.AllInOne, args)
	if err != nil {
		response.HandleError(resp, err)
		return
	}

	pipelines.NewCreateInstallerPipeline(runtime)
	response.SuccessNoData(resp)
}

func (h *Handler) handlerGreetings(req *restful.Request, resp *restful.Response) {
	logger.Infof("handler greetings req: %s", req.Request.Method)

	if err := mock.Greetings(); err != nil {
		logger.Errorf("greetings failed %v", err)
	}

	response.SuccessNoData(resp)
}

func (h *Handler) handlerDownloadEx(req *restful.Request, resp *restful.Response) {
	logger.Infof("handler download req: %s", req.Request.Method)

	arg := common.Argument{}

	if err := download.CreateDownload(arg); err != nil {
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
