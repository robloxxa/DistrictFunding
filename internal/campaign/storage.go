package campaign

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robloxxa/DistrictFunding/pkg/db"
	"time"
)

type Campaign struct {
	Id            int       `db:"id"`
	CreatorId     string    `db:"creator_id"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	Goal          int       `db:"goal"`
	CurrentAmount int       `db:"current_amount"`
	Deadline      time.Time `db:"deadline"`
	Archived      bool      `db:"archived"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type CampaignDonated struct {
	Id            int `db:"id"`
	CampaignId    int `db:"campaign_id"`
	AccountId     int `db:"account_id"`
	AmountDonated int `db:"amount_donated"`
}

type CampaignEditHistory struct {
	Id            int       `db:"id"`
	CampaignId    int       `db:"campaign_id"`
	Description   string    `db:"description"`
	Goal          int       `db:"goal"`
	CurrentAmount int       `db:"current_amount"`
	Deadline      time.Time `db:"deadline"`
	ModifiedAt    time.Time `db:"modified_at"`
}

type CampaignModel interface {
	GetById(string) (*Campaign, error)
	Create(*Campaign) error
	UpdateById(*Campaign) error
	ListByCreatorId(string) ([]Campaign, error)
}

type CampaignDonatedModel interface {
}

type CampaignEditHistoryModel interface {
}

type campaignModel struct {
	db *pgxpool.Pool
}

func (cm *campaignModel) GetById(id string) (*Campaign, error) {
	query :=
		`SELECT * FROM Campaign WHERE id = $1`

	return db.QueryOneRowToAddrStruct[Campaign](context.Background(), cm.db, query, id)
}

func (cm *campaignModel) Create(c *Campaign) error {
	sql :=
		`INSERT INTO Campaign (creator_id, name, description, goal, deadline) 
		VALUES ($1, $2, $3, $4, $5)`

	if _, err := cm.db.Exec(context.Background(), sql, c.CreatorId, c.Name, c.Description, c.Goal, c.Deadline); err != nil {
		return err
	}
	return nil
}

type campaignDonatedModel struct {
	db *pgxpool.Pool
}

type campaignEditHistoryModel struct {
	db *pgxpool.Pool
}
