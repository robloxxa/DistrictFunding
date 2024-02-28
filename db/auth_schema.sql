CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(32) NOT NULL,
    first_name TEXT,
    last_name TEXT,
    email VARCHAR(255) NOT NULL,
    -- Passwords stores as BCRYPT hash, thus the limit is 60 characters
    password VARCHAR(60) NOT NULL
    UNIQUE(username, email)
);