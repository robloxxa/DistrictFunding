package auth

import (
	"net/http"

	"github.com/go-chi/render"
)

type User struct {
	Username string `json:"username" database:"username"`
	Email string `json:"email" database:"email"`
}


// TODO: add other fields like phone number, firstname, secondname, etc.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required, email"`
	Username string `json:"login" validate:"required, gt=5, lt=20"`
	Password string `json:"password" validate:"required, gt=8"`
}

type ErrorResponse struct {
	Status int `json:"-"`
	Error string `json:"error"`
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Status)
	
	return nil
}
