package mango

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	status int
	body   any
}

func (r *jsonResponse) Send(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(r.status)
	json.NewEncoder(w).Encode(r.body)
}

func jsonRespond(status int, body any) Response {
	return &jsonResponse{body: body, status: status}
}
