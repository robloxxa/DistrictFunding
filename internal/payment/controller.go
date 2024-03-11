package payment

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robloxxa/DistrictFunding/pkg/jwtauth"
	"net/http"
)

type Api struct {
	r  chi.Router
	ja *jwtauth.JWTAuth
}

func NewController(db *pgxpool.Pool, ja *jwtauth.JWTAuth) *Api {
	a := &Api{
		chi.NewRouter(),
		ja,
	}

	a.r.Route("/campaign/{campaignId}", func(r chi.Router) {

	})
	return a
}

func (a *Api) Webhook(w http.ResponseWriter, r *http.Request) {

}

func (a *Api) DonateCampaign(w http.ResponseWriter, r *http.Request) {

}

func (a *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.r.ServeHTTP(w, r)
}
