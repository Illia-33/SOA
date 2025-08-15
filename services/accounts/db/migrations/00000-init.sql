CREATE TABLE IF NOT EXISTS accounts (
    id integer PRIMARY KEY,
    login varchar(32),
    password varchar(32),
    email varchar(320),
    phone_number varchar(18),
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE IF NOT EXISTS profiles (
    id integer PRIMARY KEY,
    account_id integer,
    name varchar(32),
    surname varchar(32),
    profile_id uuid,
    birthday date,
    bio varchar(200),
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE IF NOT EXISTS api_tokens (
    id integer PRIMARY KEY,
    account_id integer,
    token varchar(64),
    valid_untile timestamp,
    read_access boolean,
    write_access boolean,
    created_at timestamp
);
