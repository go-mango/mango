package json

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	status int
	body   any
}

func (r *JSONResponse) Send(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(r.status)
	json.NewEncoder(w).Encode(r.body)
}

func OK(body any) *JSONResponse {
	return &JSONResponse{body: body, status: http.StatusOK}
}

func Response(status int, body any) *JSONResponse {
	return &JSONResponse{body: body, status: status}
}
