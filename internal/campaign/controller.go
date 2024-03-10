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

	// TODO: list campaigns
	//a.r.Get("/", http.NotFound)

	// Campaign creating route
	a.r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(ja))
		r.Use(jwtauth.Authenticator)

		r.Post("/", a.CreateCampaign)
	})

	// All routes with campaignId as path parameter
	a.r.Route("/{campaignId}", func(r chi.Router) {
		r.Use(a.CampaignCtx)
		r.Get("/", a.GetCampaign)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(ja))
			r.Use(jwtauth.Authenticator)
			r.Use(IsCampaignOwner)

			r.Put("/", a.UpdateCampaign)
			r.Delete("/", a.DeleteCampaign)
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
			response.Error(w, http.StatusUnauthorized, err)
			return
		}

		campaign, err := CampaignFromCtx(r.Context())
		if err != nil {
			response.Error(w, http.StatusUnauthorized, err)
			return
		}

		if token.Subject() != campaign.CreatorId {
			response.Error(w, http.StatusUnauthorized, fmt.Errorf("campaign creator id is not equal to requester id"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Api) GetCampaign(w http.ResponseWriter, r *http.Request) {
	var res GetCampaignResponse

	c, err := CampaignFromCtx(r.Context())
	if err != nil {
		response.Error(w, http.StatusNotFound, err)
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

	response.Json(w, &res)
}

func (a *Api) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var (
		req CreateCampaignRequest
	)

	token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err)
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
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	res := CreateCampaignResponse{
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
	w.Header().Add("Location", fmt.Sprintf("/%d", c.Id))
	w.WriteHeader(http.StatusCreated)
	response.Json(w, &res)
}

func (a *Api) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	campaign, err := CampaignFromCtx(r.Context())
	if err != nil {
		response.Error(w, http.StatusNotFound, err)
		return
	}

	token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	if campaign.CreatorId != token.Subject() {
		response.Error(w, http.StatusBadRequest, fmt.Errorf("you can't delete campaign created by other users"))
		return
	}

	if campaign.Archived == true {
		response.Error(w, http.StatusBadRequest, fmt.Errorf("campaign already archived"))
		return
	}

	// TODO: add logic to return money back to donaters
	err = a.campaign.Archive(campaign.Id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
}

func (a *Api) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	var (
		req UpdateCampaignRequest
	)
	campaign, err := CampaignFromCtx(r.Context())
	if err != nil {
		response.Error(w, http.StatusNotFound, err)
		return
	}

	token, err := jwtauth.FromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	if campaign.CreatorId != token.Subject() {
		response.Error(w, http.StatusBadRequest, fmt.Errorf("you can't update campaign created by other users"))
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if req == (UpdateCampaignRequest{}) {
		response.Error(w, http.StatusBadRequest, fmt.Errorf("No request"))
	}

	if req.Goal != nil {
		campaign.Goal = *req.Goal
	}

	if req.Description != nil {
		campaign.Description = *req.Description
	}

	if req.Deadline != nil {
		campaign.Deadline = *req.Deadline
	}

	// TODO: make a multiple query that will insert old values from campaign
	err = a.campaign.Update(campaign)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
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
