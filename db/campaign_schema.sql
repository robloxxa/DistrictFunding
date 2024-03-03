CREATE TABLE IF NOT EXISTS "campaign" (
    id SERIAL PRIMARY KEY,
    creator_id UUID NOT NULL,
    -- TODO: think about campaign name length
    name VARCHAR(255),
    description TEXT,
    amount_needed INTEGER DEFAULT 0,
    amount_collected INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    -- TODO: make an automatic function changing updated_at
    updated_at TIMESTAMPTZ DEFAULT current_timestamp
);
