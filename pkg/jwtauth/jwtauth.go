package jwtauth

import (
	"context"
	"errors"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/robloxxa/DistrictFunding/pkg/response"
	"net/http"
	"strings"
)

var (
	TokenKey = &contextKey{"token"}
	ErrorKey = &contextKey{"tokenError"}
)

var (
	ErrTokenNotFound = errors.New("token not found")
)

type JWTAuth struct {
	jwtParser jwt.SignEncryptParseOption
}

func New(alg jwa.SignatureAlgorithm, key interface{}) *JWTAuth {
	return &JWTAuth{jwt.WithKey(alg, key)}
}

func (ja *JWTAuth) Sign(token jwt.Token) ([]byte, error) {
	return jwt.Sign(token, ja.jwtParser)
}

func (ja *JWTAuth) Parse(token string) (jwt.Token, error) {
	return jwt.Parse([]byte(token), ja.jwtParser)
}

func (ja *JWTAuth) Verify(requestParser func(r *http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := ParseTokenFromRequest(ja, r, requestParser)
			ctx = NewContext(ctx, token, err)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Verifier(ja *JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return ja.Verify(ParseTokenFromHeader)(next)
	}
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := FromContext(r.Context())
		if err != nil || token == nil {
			response.NewApiError(http.StatusUnauthorized, err).WriteResponse(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func VerifyToken(ja *JWTAuth, token string) (jwt.Token, error) {
	return jwt.Parse([]byte(token), ja.jwtParser, jwt.WithVerify(true))
}

func ParseTokenFromRequest(ja *JWTAuth, r *http.Request, tokenParser func(r *http.Request) string) (jwt.Token, error) {
	token := tokenParser(r)
	if token == "" {
		return nil, ErrTokenNotFound
	}

	return VerifyToken(ja, token)
}

func NewContext(ctx context.Context, t jwt.Token, err error) context.Context {
	ctx = context.WithValue(ctx, TokenKey, t)
	ctx = context.WithValue(ctx, ErrorKey, err)
	return ctx
}

func FromContext(ctx context.Context) (jwt.Token, error) {
	t, _ := ctx.Value(TokenKey).(jwt.Token)
	err, _ := ctx.Value(ErrorKey).(error)
	return t, err
}

func ParseTokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

type contextKey struct {
	key string
}
