package render

import "net/http"

func OK(w http.ResponseWriter, data any) error {
	return writeJSON(w, http.StatusOK, data)
}
