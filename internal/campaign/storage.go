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
	Goal          uint      `db:"goal"`
	CurrentAmount uint      `db:"current_amount"`
	Deadline      time.Time `db:"deadline"`
	Archived      bool      `db:"archived"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type CampaignDonated struct {
	Id            int  `db:"id"`
	CampaignId    int  `db:"campaign_id"`
	AccountId     int  `db:"account_id"`
	AmountDonated uint `db:"amount_donated"`
}

type CampaignEditHistory struct {
	Id            int       `db:"id"`
	CampaignId    int       `db:"campaign_id"`
	Description   string    `db:"description"`
	Goal          uint      `db:"goal"`
	CurrentAmount uint      `db:"current_amount"`
	Deadline      time.Time `db:"deadline"`
	ModifiedAt    time.Time `db:"modified_at"`
}

type CampaignModel interface {
	GetById(string) (*Campaign, error)
	Create(*Campaign) (*Campaign, error)
	Update(*Campaign) error
	Archive(int) error
	//ListByCreatorId(string) ([]Campaign, error)
}

type CampaignDonatedModel interface {
}

type CampaignEditHistoryModel interface {
	Create(*CampaignEditHistory) (*CampaignEditHistory, error)
}

type campaignModel struct {
	db *pgxpool.Pool
}

func (cm *campaignModel) GetById(id string) (*Campaign, error) {
	query :=
		`SELECT * FROM Campaign WHERE id = $1`

	return db.QueryOneRowToAddrStruct[Campaign](context.Background(), cm.db, query, id)
}

func (cm *campaignModel) Create(c *Campaign) (*Campaign, error) {
	query :=
		`INSERT INTO Campaign (creator_id, name, description, goal, deadline) 
		VALUES ($1, $2, $3, $4, $5) RETURNING *`

	return db.QueryOneRowToAddrStruct[Campaign](context.Background(), cm.db, query, c.CreatorId, c.Name, c.Description, c.Goal, c.Deadline)
}

// TODO: Maybe use map[string]interface{} instead of campaign struct?
// Update updates the campaign with new values and creates a new record in campaign history with old params
func (cm *campaignModel) Update(c *Campaign) error {
	ctx := context.Background()
	tx, err := cm.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, `INSERT INTO CampaignEditHistory (campaign_id, description, goal, deadline)
	SELECT id, description, goal, deadline FROM campaign WHERE id = $1`, c.Id); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `UPDATE Campaign SET description = $2, goal = $3, deadline = $4 WHERE id = $1`, c.Id, c.Description, c.Goal, c.Deadline); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return err
}

func (cm *campaignModel) Archive(id int) error {
	query :=
		`UPDATE Campaign SET archived = true WHERE id = $1`
	_, err := cm.db.Exec(context.Background(), query, id)
	return err
}

type campaignDonatedModel struct {
	db *pgxpool.Pool
}

type campaignEditHistoryModel struct {
	db *pgxpool.Pool
}

func (chm campaignEditHistoryModel) Create(ch *CampaignEditHistory) (*CampaignEditHistory, error) {
	query :=
		`INSERT INTO campaignedithistory (campaign_id, description, goal, deadline) VALUES ($1, $2, $3, $4)`

	return db.QueryOneRowToAddrStruct[CampaignEditHistory](context.Background(), chm.db, query, ch.CampaignId, ch.Description, ch.Goal, ch.Deadline)
}
