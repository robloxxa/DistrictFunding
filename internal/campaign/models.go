package campaign

import "time"

type GetCampaignResponse struct {
	Id              int       `json:"id"`
	CreatorId       string    `json:"creator_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	AmountNeeded    int       `json:"amount_needed"`
	AmountCollected int       `json:"amount_collected"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type (
	CreateCampaignRequest struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		AmountNeeded int    `json:"amount_needed"`
	}

	CreateCampaignResponse struct {
	}
)

type (
	UpdateCampaignRequest struct {
	}

	UpdateCampaignResponse struct {
	}
)
