package campaign

import "time"

type GetCampaignResponse struct {
	Id            int       `json:"id"`
	CreatorId     string    `json:"creator_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Goal          int       `json:"goal"`
	CurrentAmount int       `json:"current_amount"`
	Deadline      time.Time `json:"deadline"`
	Archived      bool      `json:"archived"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type (
	CreateCampaignRequest struct {
		Name          string    `json:"name"`
		Description   string    `json:"description"`
		Goal          int       `json:"goal"`
		CurrentAmount int       `json:"current_amount"`
		Deadline      time.Time `json:"deadline"`
	}

	CreateCampaignResponse struct {
		GetCampaignResponse
	}
)

type (
	UpdateCampaignRequest struct {
	}

	UpdateCampaignResponse struct {
	}
)
