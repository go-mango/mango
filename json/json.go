package json

import (
	"encoding/json"
	"net/http"

	"github.com/go-mango/mango"
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

func OK(body any) mango.Response {
	return &jsonResponse{body: body, status: http.StatusOK}
}

func Response(status int, body any) mango.Response {
	return &jsonResponse{body: body, status: status}
}
