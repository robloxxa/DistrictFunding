package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robloxxa/DistrictFunding/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	router *chi.Mux
	users  UserModel
}

func NewAuthController(db *pgxpool.Pool) *AuthController {
	c := AuthController{
		router: chi.NewRouter(),
		users:  &userModel{db},
	}

	c.router.Post("/register", c.SignUp)
	c.router.Post("/login", c.SignIn)

	c.router.Post("/logout", c.SignOut)
	c.router.Get("/me", c.Me)
	return &c
}

func (a *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.NewApiError(http.StatusBadRequest, errors.New("failed to parse json body")).Render(w, r)
		return
	}

	// Validate that request is correct
	validator := validator.New(validator.WithRequiredStructEnabled())
	if err := validator.Struct(req); err != nil {
		response.NewApiError(http.StatusBadRequest, err).Render(w, r)
		return
	}

	// Query database to see if username is already taken
	// TODO: check if empty scan can be dangerous
	if has, err := a.users.HasUsername(req.Username); has || err != nil {
		if err != nil {
			response.NewApiError(http.StatusBadRequest, err).Render(w, r)
		}
		response.NewApiError(http.StatusBadRequest, fmt.Errorf("username already exists")).Render(w, r)
		return
	}

	// Username is free, hashing password and storing user data in db
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.NewApiError(http.StatusBadRequest, err).Render(w, r)
		return
	}

	_, err = a.users.Insert(&User{Username: req.Username, Email: req.Email, FirstName: req.FirstName, LastName: req.LastName, Password: string(hash)})
	if err != nil {
		response.NewApiError(http.StatusBadRequest, err).Render(w, r)
		return
	}

	


}

func (a *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthController) SignOut(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthController) Me(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
