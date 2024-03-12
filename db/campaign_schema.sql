CREATE DATABASE CAMPAIGN_DB;

CREATE TABLE IF NOT EXISTS Campaign (
    id SERIAL PRIMARY KEY,
    creator_id UUID NOT NULL,
    -- TODO: think about campaign name length
    name VARCHAR(255),
    description TEXT,
    goal INTEGER DEFAULT 0,
    current_amount INTEGER DEFAULT 0,
    deadline TIMESTAMPTZ NOT NULL,
    archived BOOL DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    -- TODO: make an automatic function changing updated_at
    updated_at TIMESTAMPTZ DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS CampaignDonated (
    id SERIAL PRIMARY KEY,
    campaign_id INT,
    account_id UUID NOT NULL,
    amount_donated INT NOT NULL,
    CONSTRAINT fk_campaign
        FOREIGN KEY(campaign_id)
            REFERENCES Campaign(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS CampaignEditHistory (
    id SERIAL PRIMARY KEY,
    campaign_id INT NOT NULL,
    description TEXT,
    goal int,
    deadline TIMESTAMPTZ,
    modified_at timestamptz DEFAULT current_timestamp,
    CONSTRAINT fk_campaign
        FOREIGN KEY(campaign_id)
            REFERENCES Campaign(id) ON DELETE CASCADE
);
