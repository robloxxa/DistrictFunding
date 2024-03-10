package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robloxxa/DistrictFunding/pkg/jwtauth"
	"github.com/robloxxa/DistrictFunding/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

type Controller struct {
	router  *chi.Mux
	jwt     *jwtauth.JWTAuth
	account AccountModel
}

func NewController(db *pgxpool.Pool, ja *jwtauth.JWTAuth) *Controller {
	c := Controller{
		router:  chi.NewRouter(),
		account: &accountModel{db},
		jwt:     ja,
	}

	c.router.Use(jwtauth.Verifier(ja))

	c.router.Post("/signin", c.SignIn)
	c.router.Post("/signup", c.SignUp)
	c.router.Group(func(r chi.Router) {
		r.Use(jwtauth.Authenticator)
		r.Post("/signout", c.SignOut)
		r.Get("/me", c.Me)
	})

	return &c
}

func (a *Controller) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest

	token, err := jwtauth.FromContext(r.Context())
	if token != nil {
		response.Error(w, http.StatusBadRequest, errors.New("already authorized"))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, errors.New("failed to parse json body"))
		return
	}

	// Validate that request is correct
	val := validator.New(validator.WithRequiredStructEnabled())
	if err := val.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	// Query database to see if username is already taken
	// TODO: maybe make a separate route for checking username/email?
	if err := a.account.HasUsername(req.Username); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	// Username is free, hashing password and storing user data in db
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	user := &Account{Username: req.Username, Email: req.Email, FirstName: req.FirstName, LastName: req.LastName, Password: string(hash)}
	if err := a.account.Create(user); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	tokenString, err := generateJWTFromUser(a.jwt, user)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.WriteHeader(http.StatusCreated)
	response.Message(w, "User is created successfully ")
}

func (a *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	_, err := jwtauth.FromContext(r.Context())
	// TODO: see what errors could jwtauth throw in this context
	if err == nil {
		response.Error(w, http.StatusBadRequest, fmt.Errorf("already signed in"))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	user, err := a.account.FindByUsernameOrEmail(req.UsernameOrEmail)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			response.Error(w, http.StatusBadRequest, fmt.Errorf("invalid username or password"))
		default:
			response.Error(w, http.StatusBadRequest, err)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Error(w, http.StatusBadRequest, errors.New("invalid username or password"))
		return
	}

	tokenString, err := generateJWTFromUser(a.jwt, user)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Add("Authorization", "Bearer "+tokenString)
	w.WriteHeader(http.StatusOK)
}

func (a *Controller) SignOut(w http.ResponseWriter, r *http.Request) {
	// TODO: handle SignOut logic by either making blacklist in database or some other method like additional column for db or something
	_, err := jwtauth.FromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	w.Header().Del("Authorization")
}

func (a *Controller) Me(w http.ResponseWriter, r *http.Request) {
	token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	user, err := a.account.GetByUUID(token.Subject())
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	meReq := &MeRequest{
		user.Id,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
	}

	response.Json(w, meReq)
}

func (a *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func generateJWTFromUser(ja *jwtauth.JWTAuth, user *Account) (string, error) {
	token, err := jwt.NewBuilder().
		Subject(user.Id).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(24 * time.Hour)).
		Build()
	if err != nil {
		return "", err
	}
	tokenString, err := ja.Sign(token)

	return string(tokenString), err
}
