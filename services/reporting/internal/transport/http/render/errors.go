package render

import "net/http"

const (
	ErrBadRequest     = "Bad request"
	ErrInternalServer = "Internal server error"
)

func BadRequestErr(w http.ResponseWriter) {
	// TODO log error
	writeErrorJSON(w, http.StatusBadRequest, ErrBadRequest)
}

func InternalServerErr(w http.ResponseWriter) {
	writeErrorJSON(w, http.StatusInternalServerError, ErrInternalServer)
}
