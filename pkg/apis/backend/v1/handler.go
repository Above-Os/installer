package v1

import (
	"fmt"

	"bytetrade.io/web3os/installer/pkg/api/response"
	"bytetrade.io/web3os/installer/pkg/common"
	"bytetrade.io/web3os/installer/pkg/log"
	"bytetrade.io/web3os/installer/pkg/phase/download"
	"github.com/emicklei/go-restful/v3"
)

type Handler struct {
	// apis.Base
	// appService *app_service.Client
}

func New() *Handler {
	// as := app_service.NewAppServiceClient()
	// return &Handler{
	// 	appService: as,
	// }
	return &Handler{}
}

func (h *Handler) handlerTest(req *restful.Request, resp *restful.Response) {
	log.Infof("handler test req: %s", req.Request.Method)
	response.SuccessNoData(resp)
}

func (h *Handler) handlerDownload(req *restful.Request, resp *restful.Response) {
	log.Infof("handler download req: %s", req.Request.Method)

	arg := common.Argument{}

	if err := download.CreateDownload(arg, "curl -L -o %s %s"); err != nil {
		fmt.Println("---api download / err---", err)
	}

	response.SuccessNoData(resp)
}

func (h *Handler) handlerInstaller(req *restful.Request, resp *restful.Response) {
	log.Infof("handler installer req: %s", req.Request.Method)
	response.SuccessNoData(resp)
}
