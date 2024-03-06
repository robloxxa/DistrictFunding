package campaign

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robloxxa/DistrictFunding/pkg/jwtauth"
	"github.com/robloxxa/DistrictFunding/pkg/response"
	"net/http"
)

var (
	campaignKey = &contextKey{"campaign"}
	errorKey    = &contextKey{"error"}
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
		r.Post("/", a.CreateCampaign)
		r.Use(a.CampaignCtx)
		r.Get("/", a.GetCampaign)
		//r.Put("/", a.UpdateCampaign)
		r.Delete("/", a.DeleteCampaign)
	})

	return a
}

func (a *Api) CampaignCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		campaignId := chi.URLParam(r, "campaignId")
		campaign, err := a.campaign.GetById(campaignId)
		ctx := NewCampaignContext(r.Context(), campaign, err)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Api) GetCampaign(w http.ResponseWriter, r *http.Request) {
	c, err := CampaignFromCtx(r.Context())
	if err != nil {
		response.NewApiError(http.StatusNotFound, err).WriteResponse(w)
		return
	}

	res := GetCampaignResponse{
		c.Id,
		c.CreatorId,
		c.Name,
		c.Description,
		c.Goal,
		c.CurrentAmount,
		c.Deadline,
		c.Archived,
		c.CreatedAt,
		c.UpdatedAt,
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		response.NewApiError(http.StatusInternalServerError, err).WriteResponse(w)
		return
	}
}

func (a *Api) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var req CreateCampaignRequest

	if err := json.NewEncoder(w).Encode(&c); err != nil {
		response.NewApiError(http.StatusBadRequest, err).WriteResponse(w)
		return
	}

}

func (a *Api) DeleteCampaign(w http.ResponseWriter, r *http.Request) {

}

//func (a *Api) ChangeCampaign(w http.ResponseWriter, r *http.Request) {
//
//}

func (a *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.r.ServeHTTP(w, r)
}

func NewCampaignContext(ctx context.Context, c *Campaign, err error) context.Context {
	ctx = context.WithValue(ctx, campaignKey, c)
	ctx = context.WithValue(ctx, errorKey, err)
	return ctx
}

func CampaignFromCtx(ctx context.Context) (*Campaign, error) {
	campaign, _ := ctx.Value(campaignKey).(*Campaign)
	err, _ := ctx.Value(errorKey).(error)
	return campaign, err
}

type contextKey struct {
	key string
}
