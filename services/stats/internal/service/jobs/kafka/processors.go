package kafkajobs

import (
	"context"
	"soa-socialnetwork/services/stats/internal/kafka"
	"soa-socialnetwork/services/stats/internal/repo"
	"soa-socialnetwork/services/stats/pkg/models"
)

type Processor struct {
	view    topicWorker[models.PostViewEvent]
	like    topicWorker[models.PostLikeEvent]
	comment topicWorker[models.PostCommentEvent]
}

func NewProcessor(connCfg kafka.ConnectionConfig, db repo.Database) (ps Processor, err error) {
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
		return Processor{}, err
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
		return Processor{}, err
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
		return Processor{}, err
	}

	return
}

func (r *Processor) Start(ctx context.Context) {
	r.view.start(ctx)
	r.like.start(ctx)
	r.comment.start(ctx)
}

func (r *Processor) Close() {
	r.view.close()
	r.like.close()
	r.comment.close()
}
