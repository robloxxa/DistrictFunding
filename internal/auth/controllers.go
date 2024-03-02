package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robloxxa/DistrictFunding/pkg/jwtauth"
	"github.com/robloxxa/DistrictFunding/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

type Api struct {
	router *chi.Mux
	jwt    *jwtauth.JWTAuth
	users  UserModel
}

func NewAuthController(db *pgxpool.Pool, ja *jwtauth.JWTAuth) *Api {
	c := Api{
		router: chi.NewRouter(),
		users:  &userModel{db},
		jwt:    ja,
	}

	c.router.Use(jwtauth.Verifier(ja))
	c.router.Use(jwtauth.Authenticator)

	c.router.Post("/signup", c.SignUp)
	c.router.Post("/signin", c.SignIn)
	c.router.Post("/signout", c.SignOut)

	c.router.Get("/me", c.Me)
	return &c
}

func (a *Api) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.NewApiError(http.StatusBadRequest, errors.New("failed to parse json body")).Render(w, r)
		return
	}

	// Validate that request is correct
	val := validator.New(validator.WithRequiredStructEnabled())
	if err := val.Struct(req); err != nil {
		response.NewApiError(http.StatusBadRequest, err).Render(w, r)
		return
	}

	// Query database to see if username is already taken
	// TODO: maybe make a separate route for checking username/email?
	if err := a.users.HasUsername(req.Username); err != nil {
		response.NewApiError(http.StatusBadRequest, err).Render(w, r)
		return
	}

	// Username is free, hashing password and storing user data in db
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.NewApiError(http.StatusBadRequest, err).Render(w, r)
		return
	}
	user := &User{Username: req.Username, Email: req.Email, FirstName: req.FirstName, LastName: req.LastName, Password: string(hash)}
	if err := a.users.Insert(user); err != nil {
		response.NewApiError(http.StatusBadRequest, err).Render(w, r)
		return
	}
	tokenString, err := generateJWTFromUser(a.jwt, user)
	if err != nil {
		response.NewApiError(http.StatusBadRequest, err).Render(w, r)
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)
}

func (a *Api) SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	_, _, err := jwtauth.FromContext(r.Context())
	// TODO: see what errors could jwtauth throw in this context
	if err == nil {
		response.NewApiError(http.StatusBadRequest, fmt.Errorf("already signed in")).Render(w, r)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NewApiError(http.StatusBadRequest, err).Render(w, r)
		return
	}
	user, err := a.users.FindByUsernameOrEmail(req.UsernameOrEmail)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			response.NewApiError(http.StatusBadRequest, fmt.Errorf("invalid username or password")).Render(w, r)
		default:
			response.NewApiError(http.StatusBadRequest, err)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.NewApiError(http.StatusBadRequest, errors.New("invalid username or password")).Render(w, r)
		return
	}

	tokenString, err := generateJWTFromUser(a.jwt, user)
	if err != nil {
		response.NewApiError(http.StatusBadRequest, err).Render(w, r)
		return
	}

	w.Header().Add("Authorization", tokenString)
}

func (a *Api) SignOut(w http.ResponseWriter, r *http.Request) {
	// TODO: handle SignOut logic by either making blacklist in database or some other method like additional column for db or something
	_, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		response.NewApiError(http.StatusUnauthorized, err).Render(w, r)
		return
	}

	w.Header().Del("Authorization")
}

func (a *Api) Me(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		response.NewApiError(http.StatusUnauthorized, err).Render(w, r)
		return
	}

	user, err := a.users.GetByUUID(claims["sub"].(string))
	if err != nil {
		response.NewApiError(http.StatusBadRequest, err).Render(w, r)
		return
	}

	meReq := &MeRequest{
		user.Id,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
	}

	if err := json.NewEncoder(w).Encode(&meReq); err != nil {
		response.NewApiError(http.StatusBadRequest, err).Render(w, r)
		return
	}
}

func (a *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func generateJWTFromUser(ja *jwtauth.JWTAuth, user *User) (string, error) {
	claims := map[string]interface{}{
		"sub":      user.Id,
		"issuedAt": time.Now(),
		"exp":      time.Now().Add(time.Hour * 24),
	}
	_, tokenString, err := ja.Encode(claims)

	return tokenString, err
}
