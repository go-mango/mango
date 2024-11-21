package mango

import (
	"net/http"
)

type Response interface {
	Send(http.ResponseWriter)
}

type responseWithStatus struct {
	ctx *Context
	w   http.ResponseWriter
}

func (r *responseWithStatus) Header() http.Header {
	return r.w.Header()
}

func (r *responseWithStatus) Write(p []byte) (int, error) {
	return r.w.Write(p)
}

func (r *responseWithStatus) WriteHeader(status int) {
	r.ctx.status = status
	r.w.WriteHeader(status)
}

type textResponse struct {
	status int
	body   string
}

func (r *textResponse) Send(w http.ResponseWriter) {
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(r.status)
	w.Write([]byte(r.body))
}

func text(status int, body string) Response {
	return &textResponse{body: body, status: status}
}
