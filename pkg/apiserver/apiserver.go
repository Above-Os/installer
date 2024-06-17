package apiserver

import (
	"net/http"

	"bytetrade.io/web3os/installer/pkg/api/response"
	apisV1alpha1 "bytetrade.io/web3os/installer/pkg/apis/backend/v1"
	"bytetrade.io/web3os/installer/pkg/constants"
	"bytetrade.io/web3os/installer/pkg/log"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	urlruntime "k8s.io/apimachinery/pkg/util/runtime"
)

type APIServer struct {
	Server    *http.Server
	container *restful.Container
}

func New() (*APIServer, error) {
	s := &APIServer{}

	server := &http.Server{
		Addr: constants.ApiServerListenAddress,
	}

	s.Server = server
	return s, nil
}

func (s *APIServer) PrepareRun() error {
	s.container = restful.NewContainer()
	s.container.RecoverHandler(logStackOnRecover)
	s.container.Filter(func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		defer func() {
			if e := recover(); e != nil {
				response.HandleInternalError(resp, errors.Errorf("server internal error: %v", e))
			}
		}()

		chain.ProcessFilter(req, resp)
	})
	// s.container.Filter(authenticate)
	s.container.Router(restful.CurlyRouter{})

	s.installStaticResources()
	s.installModuleAPI()
	s.installAPIDocs()

	var modulePaths []string
	for _, ws := range s.container.RegisteredWebServices() {
		modulePaths = append(modulePaths, ws.RootPath())
	}
	log.Infow("registered module", "paths", modulePaths)

	s.Server.Handler = s.container
	return nil
}

func (s *APIServer) Run() error {
	err := s.Server.ListenAndServe()
	if err != nil {
		return errors.Errorf("listen and serve err: %v", err)
	}
	return nil
}

func (s *APIServer) installAPIDocs() {
	config := restfulspec.Config{
		WebServices:                   s.container.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "./apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	s.container.Add(restfulspec.NewOpenAPIService(config))
}

func (s *APIServer) installStaticResources() {
	ws := &restful.WebService{}

	ws.Route(ws.GET("/web/{subpath:*}").To(staticFromPathParam))
	// ws.Route(ws.GET("/web").To(staticFromQueryParam))

	s.container.Add(ws)
}

func (s *APIServer) installModuleAPI() {
	urlruntime.Must(apisV1alpha1.AddContainer(s.container))
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Installer API Server Docs",
			Description: "Backend For Installer",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "bytetrade",
					Email: "dev@bytetrade.io",
					URL:   "http://bytetrade.io",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "Apache License 2.0",
					URL:  "http://www.apache.org/licenses/LICENSE-2.0",
				},
			},
			Version: "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{{TagProps: spec.TagProps{
		Name:        "Installer",
		Description: "Terminus Installer"}}}
}
