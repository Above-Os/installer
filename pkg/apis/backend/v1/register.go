package v1

import (
	"net/http"

	"bytetrade.io/web3os/installer/pkg/api/response"
	"bytetrade.io/web3os/installer/pkg/apiserver/runtime"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
)

var ModuleVersion = runtime.ModuleVersion{Name: "backend", Version: "v1"}

var tags = []string{"backend"}

func AddContainer(c *restful.Container) error {
	ws := runtime.NewWebService(ModuleVersion)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)

	handler := New()

	ws.Route(ws.POST("/test").
		To(handler.handlerTest).
		Doc("").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(http.StatusOK, "", response.Response{}))

	ws.Route(ws.POST("/download").
		To(handler.handlerDownload).
		Doc("").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(http.StatusOK, "", response.Response{}))

	ws.Route(ws.POST("/install").
		To(handler.handlerInstallKk).
		Doc("").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(http.StatusOK, "", response.Response{}))

	ws.Route(ws.POST("/install_terminus").
		To(handler.handlerInstallTerminus).
		Doc("").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(http.StatusOK, "", response.Response{}))

	c.Add(ws)

	return nil

}
