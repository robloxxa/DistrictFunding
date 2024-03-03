package campaign

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robloxxa/DistrictFunding/pkg/jwtauth"
	"github.com/robloxxa/DistrictFunding/pkg/response"
	"net/http"
)

var (
	campaignKey = &contextKey{"campaign"}
)

type Api struct {
	r        chi.Router
	ja       *jwtauth.JWTAuth
	campaign CampaignModel
}

func NewController(db *pgxpool.Pool, ja *jwtauth.JWTAuth) *Api {
	a := &Api{
		chi.NewRouter(),
		ja,
		&campaignModel{db},
	}

	// TODO:
	a.r.Get("/", http.NotFound)

	a.r.Route("/{campaignId}", func(r chi.Router) {
		r.Use(a.CampaignCtx)
		r.Get("/", a.GetCampaign)
		//r.Post("/", a.UpdateCampaign)
		r.Put("/", a.CreateCampaign)
		r.Delete("/", a.DeleteCampaign)
	})

	return a
}

func (a *Api) CampaignCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		campaignId := chi.URLParam(r, "campaignId")
		campaign, err := a.campaign.Get(campaignId)
		if err != nil {
			response.NewApiError(http.StatusNotFound, fmt.Errorf("couldn't found campaign id: %w", err))
			return
		}
		ctx := context.WithValue(r.Context(), campaignKey, campaign)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Api) GetCampaign(w http.ResponseWriter, r *http.Request) {
	campaign, err := CampaignFromCtx(r.Context())
	if err != nil {
		response.NewApiError(http.StatusNotFound, err)
	}

}

func (a *Api) CreateCampaign(w http.ResponseWriter, r *http.Request) {

}

func (a *Api) DeleteCampaign(w http.ResponseWriter, r *http.Request) {

}

//func (a *Api) ChangeCampaign(w http.ResponseWriter, r *http.Request) {
//
//}

func (a *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.r.ServeHTTP(w, r)
}

func CampaignFromCtx(ctx context.Context) (*Campaign, error) {
	campaign, ok := ctx.Value(campaignKey).(*Campaign)
	if !ok {
		return nil, fmt.Errorf("no campaign found")
	}
	return campaign, nil
}

type contextKey struct {
	key string
}
