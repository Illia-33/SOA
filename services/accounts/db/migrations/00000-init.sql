CREATE TABLE IF NOT EXISTS accounts (
    id INTEGER PRIMARY KEY,
    login VARCHAR(32),
    password VARCHAR(32),
    email VARCHAR(320),
    phone_number VARCHAR(18),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS profiles (
    id INTEGER PRIMARY KEY,
    account_id INTEGER,
    name VARCHAR(32),
    surname VARCHAR(32),
    profile_id UUID,
    birthday DATE,
    bio VARCHAR(200),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS api_tokens (
    id INTEGER PRIMARY KEY,
    account_id INTEGER,
    token VARCHAR(64),
    valid_until TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    read_access BOOLEAN,
    write_access BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
