package campaign

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Campaign struct {
	Id              int       `database:"id"`
	CreatorId       string    `database:"creator_id"`
	Name            string    `database:"name"`
	Description     string    `database:"description"`
	AmountNeeded    int       `database:"amount_needed"`
	AmountCollected int       `database:"amount_collected"`
	CreatedAt       time.Time `database:"created_at"`
	UpdatedAt       time.Time `database:"updated_at"`
}

type CampaignModel interface {
	Get(string) (*Campaign, error)
	Insert(*Campaign) error
	UpdateById(string, *Campaign) error
	ListByCreatorId(string) ([]Campaign, error)
}

type campaignModel struct {
	db *pgxpool.Pool
}

func (cm *campaignModel) Get(id string) (*Campaign, error) {
	var c Campaign

	query :=
		`SELECT (id, creator_id, name, description, amount_needed, amount_collected) FROM "campaigns" WHERE id = $1`

	if err := cm.db.QueryRow(context.Background(), query, id).Scan(&c); err == nil {
		return nil, err
	}

	return &c, nil
}
