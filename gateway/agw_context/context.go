package agw_context

import "net/http"

type AGWContext struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	Response *http.Response
}

func NewAGWContext(w http.ResponseWriter, r *http.Request) *AGWContext {
	return &AGWContext{
		Request: r,
		Writer:  w,
	}
}
