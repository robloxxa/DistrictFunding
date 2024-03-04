package campaign

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Campaign struct {
	Id            int       `database:"id"`
	CreatorId     string    `database:"creator_id"`
	Name          string    `database:"name"`
	Description   string    `database:"description"`
	Goal          int       `database:"goal"`
	CurrentAmount int       `database:"current_amount"`
	Deadline      time.Time `database:"deadline"`
	Archived      bool      `database:"archived"`
	CreatedAt     time.Time `database:"created_at"`
	UpdatedAt     time.Time `database:"updated_at"`
}

type CampaignDonated struct {
	Id         int `database:"id"`
	CampaignId int `database:"campaign_id"`
}

type CampaignModel interface {
	GetById(string) (*Campaign, error)
	Create(*Campaign) error
	UpdateById(*Campaign) error
	ListByCreatorId(string) ([]Campaign, error)
}

type campaignModel struct {
	db *pgxpool.Pool
}

func (cm *campaignModel) GetById(id string) (*Campaign, error) {
	var c Campaign

	query :=
		`SELECT (id, creator_id, name, description, amount_needed, amount_collected) FROM Campaign WHERE id = $1`

	if err := cm.db.QueryRow(context.Background(), query, id).Scan(&c); err == nil {
		return nil, err
	}

	return &c, nil
}

func (cm *campaignModel) Create(c *Campaign) error {
	query :=
		`INSERT INTO Campaign (creator_id, name, description, goal, current_amount, deadline) VALUES ()`

	if err := cm.db.QueryRow(context.Background(), query, id).Scan(&c); err == nil {
		return nil, err
	}

	return &c, nil
}

func (cm *campaignModel) GetById(id string) (*Campaign, error) {
	var c Campaign

	query :=
		`SELECT (id, creator_id, name, description, amount_needed, amount_collected) FROM Campaign WHERE id = $1`

	if err := cm.db.QueryRow(context.Background(), query, id).Scan(&c); err == nil {
		return nil, err
	}

	return &c, nil
}
