package shared

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
)



type ErrorResponse struct {
	status int `json:"-"`
	Err error `json:"error"`
}

func NewResponseError(status int, err error) *ErrorResponse {
	return &ErrorResponse{
		status, 
		err,
	}
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.status)
	return json.NewEncoder(w).Encode(e)
}

func (e *ErrorResponse) Error() string {
	return e.Err.Error()
}