CREATE DATABASE PAYMENT_DB;

CREATE TABLE IF NOT EXISTS Payment (
    id SERIAL PRIMARY KEY,
    payment_id VARCHAR(36) UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    campaign_id int NOT NULL,
    amount float NOT NULL,
    currency VARCHAR(3) DEFAULT 'RUB',
    returned_at timestamptz,
    created_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS Payout (
    id SERIAL PRIMARY KEY,
    payout_id VARCHAR(36) UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    campaign_id INT NOT NULL,
    amount float NOT NULL,
    currency VARCHAR(3) DEFAULT 'RUB',
    created_at TIMESTAMPTZ default current_timestamp
);