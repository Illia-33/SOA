package kafkajobs

import (
	"context"
	"soa-socialnetwork/services/stats/internal/kafka"
	"soa-socialnetwork/services/stats/internal/repo"
	"soa-socialnetwork/services/stats/pkg/models"
)

type processors struct {
	view    topicWorker[models.PostViewEvent]
	like    topicWorker[models.PostLikeEvent]
	comment topicWorker[models.PostCommentEvent]
}

func NewProcessors(connCfg kafka.ConnectionConfig, db repo.Database) (ps processors, err error) {
	defer func() {
		if err != nil {
			ps.Close()
		}
	}()

	ps.view, err = newTopicProcessor(
		connCfg,
		kafka.ConsumerConfig{
			Topic:   "view",
			GroupId: "stats-service-view",
		},
		func(ctx context.Context, batch messageBatch[models.PostViewEvent]) error {
			events := make([]models.PostViewEvent, len(batch))
			for i := range batch {
				events[i] = batch[i].Value
			}
			return db.PostsViews(ctx).Put(events...)
		},
	)
	if err != nil {
		return processors{}, err
	}

	ps.like, err = newTopicProcessor(
		connCfg,
		kafka.ConsumerConfig{
			Topic:   "like",
			GroupId: "stats-service-like",
		},
		func(ctx context.Context, batch messageBatch[models.PostLikeEvent]) error {
			events := make([]models.PostLikeEvent, len(batch))
			for i := range batch {
				events[i] = batch[i].Value
			}
			return db.PostsLikes(ctx).Put(events...)
		},
	)
	if err != nil {
		return processors{}, err
	}

	ps.comment, err = newTopicProcessor(
		connCfg,
		kafka.ConsumerConfig{
			Topic:   "comment",
			GroupId: "stats-service-comment",
		},
		func(ctx context.Context, batch messageBatch[models.PostCommentEvent]) error {
			events := make([]models.PostCommentEvent, len(batch))
			for i := range batch {
				events[i] = batch[i].Value
			}
			return db.PostsComments(ctx).Put(events...)
		},
	)
	if err != nil {
		return processors{}, err
	}

	return
}

func (r *processors) Start(ctx context.Context) {
	r.view.start(ctx)
	r.like.start(ctx)
	r.comment.start(ctx)
}

func (r *processors) Close() {
	r.view.close()
	r.like.close()
	r.comment.close()
}
