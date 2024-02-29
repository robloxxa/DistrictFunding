package response

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"

	"github.com/go-chi/render"
)



type ApiError struct {
	status int `json:"-"`
	ErrorMsg string `json:"error"`
}

func NewApiError(status int, err error) *ApiError {
	_, file, line, _ := runtime.Caller(1)
	log.Println(file, line, err)
	return &ApiError{
		status, 
		err.Error(),
	}
}

func (e *ApiError) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.status)
	return json.NewEncoder(w).Encode(e)
}
