package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robloxxa/DistrictFunding/shared"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	router *chi.Mux
	db     *pgxpool.Pool
}

func NewAuthController(db *pgxpool.Pool) *AuthController {
	c := AuthController{
		chi.NewRouter(),
		db,
	}

	c.router.Post("/register", c.Register)
	c.router.Post("/login", c.Login)

	c.router.Post("/logout", c.Logout)
	c.router.Get("/me", c.Me)
	return &c
}


func (a *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		shared.NewResponseError(http.StatusBadRequest, err).Render(w, r)
		return
	}

	// Validate that request is correct
	validator := validator.New(validator.WithRequiredStructEnabled())
	if err := validator.Struct(req); err != nil {
		shared.NewResponseError(http.StatusBadRequest, err).Render(w, r)
		return
	}

	// Query database to see if username is already taken
	// TODO: check if empty scan can be dangerous
	if err := a.db.QueryRow(context.Background(), "SELECT username FROM users WHERE username = $1", req.Username).Scan(); err != nil {
		shared.NewResponseError(http.StatusBadRequest, err).Render(w, r)
		return
	}

	// Username is free, hashing password and storing data in db
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		shared.NewResponseError(http.StatusBadRequest, err).Render(w, r)
		return
	}

	_, err = a.db.Exec(context.Background(), "INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", req.Username



}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthController) Logout(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthController) Me(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}


