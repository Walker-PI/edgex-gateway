package agw_context

import (
	"context"
	"net/http"
)

type AGWContext struct {
	ResponseWriter http.ResponseWriter
	OriginRequest  *http.Request
	ForwardRequest *http.Request
}

func NewAGWContext(writer http.ResponseWriter, req *http.Request) *AGWContext {
	return &AGWContext{
		ResponseWriter: writer,
		OriginRequest:  req,
		ForwardRequest: copyRequest(req),
	}
}

func copyRequest(req *http.Request) *http.Request {
	return req.Clone(context.Background())
}
