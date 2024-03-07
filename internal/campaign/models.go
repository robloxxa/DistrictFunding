package campaign

import "time"

type GetCampaignResponse struct {
	Id            int       `json:"id"`
	CreatorId     string    `json:"creator_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Goal          uint      `json:"goal"`
	CurrentAmount uint      `json:"current_amount"`
	Deadline      time.Time `json:"deadline"`
	Archived      bool      `json:"archived"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type (
	// TODO: add verify to fields
	CreateCampaignRequest struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Goal        uint      `json:"goal"`
		Deadline    time.Time `json:"deadline"`
	}

	//CreateCampaignResponse struct {
	//	GetCampaignResponse
	//}
)

type (
	UpdateCampaignRequest struct {
		Description *string    `json:"description,omitempty"`
		Goal        *uint      `json:"goal,omitempty"`
		Deadline    *time.Time `json:"deadline,omitempty"`
	}

	UpdateCampaignResponse struct {
	}
)
