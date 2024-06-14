package apiserver

import (
	"bytes"
	"io"
	"net/http"
	"path"

	"bytetrade.io/web3os/installer/frontend"
	"bytetrade.io/web3os/installer/pkg/log"
	restful "github.com/emicklei/go-restful/v3"
)

func staticFromPathParam(req *restful.Request, resp *restful.Response) {
	log.Infof("handler static req: %s", req.Request.Method)
	subpath := req.PathParameter("subpath")
	actual := path.Join("dist", subpath)

	file, err := frontend.Assets().Open(actual)
	if err != nil {
		http.NotFound(resp.ResponseWriter, req.Request)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		http.NotFound(resp.ResponseWriter, req.Request)
		return
	}

	if fileInfo.IsDir() {
		indexFilePath := path.Join(actual, "index.html")
		indexFile, err := frontend.Assets().Open(indexFilePath)
		if err != nil {
			http.NotFound(resp.ResponseWriter, req.Request)
			return
		}
		defer indexFile.Close()

		content, err := io.ReadAll(indexFile)
		if err != nil {
			http.NotFound(resp.ResponseWriter, req.Request)
			return
		}

		reader := bytes.NewReader(content)

		http.ServeContent(resp.ResponseWriter, req.Request, indexFilePath, fileInfo.ModTime(), reader)
	} else {

		content, err := io.ReadAll(file)
		if err != nil {
			http.NotFound(resp.ResponseWriter, req.Request)
			return
		}

		reader := bytes.NewReader(content)

		http.ServeContent(resp.ResponseWriter, req.Request, actual, fileInfo.ModTime(), reader)
	}
}

// ! ⚠️ 暂时废弃
func staticFromQueryParam(req *restful.Request, resp *restful.Response) {
	resource := req.QueryParameter("resource") // todo why here is resource?
	actual := path.Join("dist", resource)

	file, err := frontend.Assets().Open(actual)
	if err != nil {
		http.NotFound(resp.ResponseWriter, req.Request)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		http.NotFound(resp.ResponseWriter, req.Request)
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		http.NotFound(resp.ResponseWriter, req.Request)
		return
	}

	reader := bytes.NewReader(content)

	http.ServeContent(resp.ResponseWriter, req.Request, actual, fileInfo.ModTime(), reader)
}
