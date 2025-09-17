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

-- agg_post_metrics

CREATE TABLE IF NOT EXISTS agg_post_metrics(
    post_id Int32,
    metric LowCardinality(String),
    cnt AggregateFunction(count, UInt64)
)
ENGINE = AggregatingMergeTree
ORDER BY (metric, post_id);

CREATE MATERIALIZED VIEW mv_agg_post_metrics_view_count
TO agg_post_metrics
AS
SELECT
    post_id,
    'view_count' AS metric,
    countState() AS cnt
FROM posts_views
GROUP BY post_id;

CREATE MATERIALIZED VIEW mv_agg_post_metrics_like_count
TO agg_post_metrics
AS
SELECT
    post_id,
    'like_count' AS metric,
    countState() AS cnt
FROM posts_likes
GROUP BY post_id;

CREATE MATERIALIZED VIEW mv_agg_post_metrics_comment_count
TO agg_post_metrics
AS
SELECT
    post_id,
    'comment_count' AS metric,
    countState() AS cnt
FROM posts_comments
GROUP BY post_id;

-- agg_user_metrics

CREATE TABLE IF NOT EXISTS agg_user_metrics(
    account_id Int32,
    metric LowCardinality(String),
    cnt AggregateFunction(count, UInt64)
)
ENGINE = AggregatingMergeTree
ORDER BY (metric, account_id);

CREATE MATERIALIZED VIEW mv_agg_user_metrics_view_count
TO agg_user_metrics
AS
SELECT
    p.author_id AS account_id,
    'view_count' AS metric,
    countState() AS cnt
FROM posts_views pv
INNER JOIN posts p ON pv.post_id = p.post_id 
GROUP BY p.author_id;

CREATE MATERIALIZED VIEW mv_agg_user_metrics_like_count
TO agg_user_metrics
AS
SELECT
    p.author_id AS account_id,
    'like_count' AS metric,
    countState() AS cnt
FROM posts_likes pl
INNER JOIN posts p ON pl.post_id = p.post_id 
GROUP BY p.author_id;

CREATE MATERIALIZED VIEW mv_agg_user_metrics_comment_count
TO agg_user_metrics
AS
SELECT
    p.author_id AS account_id,
    'comment_count' AS metric,
    countState() AS cnt
FROM posts_comments pc
INNER JOIN posts p ON pc.post_id = p.post_id 
GROUP BY p.author_id;
