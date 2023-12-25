package rest

import (
	"encoding/json"
	"io"
	"net/http"
)

func ReadJSON[T any](closer io.ReadCloser, out *T) error {
	return json.NewDecoder(closer).Decode(&out)
}

func ServeJSON[T any](status int, w http.ResponseWriter, payload T) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}
