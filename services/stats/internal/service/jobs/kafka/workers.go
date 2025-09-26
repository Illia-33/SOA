package kafkajobs

import (
	"context"
	"soa-socialnetwork/services/stats/internal/kafka"
	"soa-socialnetwork/services/stats/internal/repo"
	"soa-socialnetwork/services/stats/pkg/models"
)

type Workers struct {
	view         topicWorker[models.PostViewEvent]
	like         topicWorker[models.PostLikeEvent]
	comment      topicWorker[models.PostCommentEvent]
	registration topicWorker[models.RegistrationEvent]
	posts        topicWorker[models.PostEvent]
}

func NewWorkers(connCfg kafka.ConnectionConfig, db repo.Database) (w Workers, err error) {
	defer func() {
		if err != nil {
			w.Close()
		}
	}()

	w.view, err = newTopicWorker(
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
		return Workers{}, err
	}

	w.like, err = newTopicWorker(
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
		return Workers{}, err
	}

	w.comment, err = newTopicWorker(
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
		return Workers{}, err
	}

	w.registration, err = newTopicWorker(
		connCfg,
		kafka.ConsumerConfig{
			Topic:   "registration",
			GroupId: "stats-service-registration",
		},
		func(ctx context.Context, batch messageBatch[models.RegistrationEvent]) error {
			events := make([]models.RegistrationEvent, len(batch))
			for i := range batch {
				events[i] = batch[i].Value
			}
			return db.Registrations(ctx).Put(events...)
		},
	)
	if err != nil {
		return Workers{}, err
	}

	w.posts, err = newTopicWorker(
		connCfg,
		kafka.ConsumerConfig{
			Topic:   "post",
			GroupId: "stats-service-post",
		},
		func(ctx context.Context, batch messageBatch[models.PostEvent]) error {
			events := make([]models.PostEvent, len(batch))
			for i := range batch {
				events[i] = batch[i].Value
			}
			return db.Posts(ctx).Put(events...)
		},
	)
	if err != nil {
		return Workers{}, err
	}

	return
}

func (r *Workers) Start(ctx context.Context) {
	r.view.start(ctx)
	r.like.start(ctx)
	r.comment.start(ctx)
	r.registration.start(ctx)
	r.posts.start(ctx)
}

func (r *Workers) Close() {
	r.view.close()
	r.like.close()
	r.comment.close()
	r.registration.close()
	r.posts.close()
}
