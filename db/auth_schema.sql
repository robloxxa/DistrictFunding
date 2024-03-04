CREATE TABLE IF NOT EXISTS Account (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(32) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    first_name TEXT,
    last_name TEXT,
    -- Passwords stores as BCRYPT hash, thus the limit is 60 characters
    password VARCHAR(60) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    -- TODO: make an automatic function changing updated_at
    updated_at TIMESTAMPTZ DEFAULT current_timestamp
);
