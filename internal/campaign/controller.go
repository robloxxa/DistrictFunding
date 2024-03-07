package campaign

import (
	"context"
	"encoding/json"
	"fmt"
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
	r               chi.Router
	ja              *jwtauth.JWTAuth
	campaign        CampaignModel
	campaignHistory CampaignEditHistoryModel
	campaignDonated CampaignDonatedModel
}

func NewController(db *pgxpool.Pool, ja *jwtauth.JWTAuth) *Api {
	a := &Api{
		chi.NewRouter(),
		ja,
		&campaignModel{db},
		&campaignEditHistoryModel{db},
		&campaignDonatedModel{db},
	}

	// TODO:

	a.r.Get("/", http.NotFound)
	a.r.Route("/{campaignId}", func(r chi.Router) {
		r.Use(a.CampaignCtx)
		r.Get("/", a.GetCampaign)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(ja))
			r.Use(jwtauth.Authenticator)

			r.Post("/", a.CreateCampaign)
			r.With(IsCampaignOwner).Put("/", a.UpdateCampaign)
			r.With(IsCampaignOwner).Delete("/", a.DeleteCampaign)
		})
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

func IsCampaignOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwtauth.FromContext(r.Context())
		if err != nil {
			response.NewApiError(http.StatusUnauthorized, err).WriteResponse(w)
			return
		}

		campaign, err := CampaignFromCtx(r.Context())
		if err != nil {
			response.NewApiError(http.StatusUnauthorized, err).WriteResponse(w)
			return
		}

		if token.Subject() != campaign.CreatorId {
			response.NewApiError(http.StatusUnauthorized, fmt.Errorf("campaign creator id is not equal to requester id")).WriteResponse(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Api) GetCampaign(w http.ResponseWriter, r *http.Request) {
	var res GetCampaignResponse

	c, err := CampaignFromCtx(r.Context())
	if err != nil {
		response.NewApiError(http.StatusNotFound, err).WriteResponse(w)
		return
	}

	res = GetCampaignResponse{
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
	var (
		req CreateCampaignRequest
	)

	token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		response.NewApiError(http.StatusUnauthorized, err).WriteResponse(w)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NewApiError(http.StatusBadRequest, err).WriteResponse(w)
		return
	}

	c := &Campaign{
		CreatorId:   token.Subject(),
		Name:        req.Name,
		Description: req.Description,
		Goal:        req.Goal,
		Deadline:    req.Deadline,
	}

	c, err = a.campaign.Create(c)
	if err != nil {
		response.NewApiError(http.StatusBadRequest, err).WriteResponse(w)
		return
	}
}

func (a *Api) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	campaign, err := CampaignFromCtx(r.Context())
	if err != nil {
		response.NewApiError(http.StatusNotFound, err).WriteResponse(w)
		return
	}

	token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		response.NewApiError(http.StatusUnauthorized, err).WriteResponse(w)
		return
	}

	if campaign.CreatorId != token.Subject() {
		response.NewApiError(http.StatusBadRequest, fmt.Errorf("you can't delete campaign created by other users")).WriteResponse(w)
		return
	}

	if campaign.Archived == true {
		response.NewApiError(http.StatusBadRequest, fmt.Errorf("campaign already archived")).WriteResponse(w)
		return
	}

	// TODO: add logic to return money back to donaters
	err = a.campaign.Archive(campaign.Id)
	if err != nil {
		response.NewApiError(http.StatusBadRequest, err).WriteResponse(w)
		return
	}
}

func (a *Api) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	var (
		req             UpdateCampaignRequest
		campaignHistory CampaignEditHistory
	)
	campaign, err := CampaignFromCtx(r.Context())
	if err != nil {
		response.NewApiError(http.StatusNotFound, err).WriteResponse(w)
		return
	}

	token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		response.NewApiError(http.StatusUnauthorized, err).WriteResponse(w)
		return
	}

	if campaign.CreatorId != token.Subject() {
		response.NewApiError(http.StatusBadRequest, fmt.Errorf("you can't update campaign created by other users")).WriteResponse(w)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NewApiError(http.StatusBadRequest, err).WriteResponse(w)
		return
	}

	if req.Goal != nil {
		campaignHistory.Goal = campaign.Goal
		campaign.Goal = *req.Goal
	}
	if req.Description != nil {
		campaignHistory.Description = campaign.Description
		campaign.Description = *req.Description
	}
	if req.Deadline != nil {
		campaignHistory.Deadline = campaign.Deadline
		campaign.Deadline = *req.Deadline
	}

	err = a.campaign.Update(campaign)
	if err != nil {
		response.NewApiError(http.StatusBadRequest, err).WriteResponse(w)
		return
	}
	a.campaignHistory.Create()

}

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
