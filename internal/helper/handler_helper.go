package helper

import (
	"encoding/json"
	"net/http"
)

func DecodeJSONBody[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var req T
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		var zero T
		return zero, err
	}
	return req, nil
}
