ALTER TABLE posts ADD COLUMN views_count INTEGER NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS likes (
    post_id INTEGER NOT NULL,
    author_account_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (post_id, author_account_id)
);

