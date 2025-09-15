CREATE TABLE IF NOT EXISTS posts_views(
    post_id Int32,
    viewer_account_id Int32,
    view_time DateTime
)
ENGINE = MergeTree
ORDER BY (post_id, view_time)
PARTITION BY toYYYYMM(view_time);

CREATE TABLE IF NOT EXISTS posts_likes(
    post_id Int32,
    liker_account_id Int32,
    like_time DateTime
)
ENGINE = MergeTree
ORDER BY (post_id, like_time)
PARTITION BY toYYYYMM(like_time);

CREATE TABLE IF NOT EXISTS posts_comments(
    post_id Int32,
    author_account_id Int32,
    comment_id Int32,
    post_time DateTime
)
ENGINE = MergeTree
ORDER BY (post_id, post_time)
PARTITION BY toYYYYMM(post_time);

CREATE TABLE IF NOT EXISTS registrations(
    account_id Int32,
    profile_id String,
    register_time DateTime
)
ENGINE = MergeTree
ORDER BY (account_id, register_time)
PARTITION BY toYYYYMM(register_time);

CREATE TABLE IF NOT EXISTS posts(
    post_id Int32,
    author_id Int32,
    post_time DateTime
)
ENGINE = MergeTree
ORDER BY (author_id, post_time)
PARTITION BY toYYYYMM(post_time);
