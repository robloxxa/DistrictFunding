package response

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/render"
)



type ApiError struct {
	status int `json:"-"`
	ErrorMsg string `json:"error"`
}

func NewApiError(status int, err error) *ApiError {
	log.Println(err)
	return &ApiError{
		status, 
		err.Error(),
	}
}

func (e *ApiError) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.status)
	return json.NewEncoder(w).Encode(e)
}
