package kafkajobs

import (
	"context"
	"errors"
	"soa-socialnetwork/services/stats/internal/kafka"
	"soa-socialnetwork/services/stats/internal/repo"
	"soa-socialnetwork/services/stats/pkg/models"
)

type Workers struct {
	workers []iTopicWorker
}

func NewWorkers(connCfg kafka.ConnectionConfig, db repo.Database) (Workers, error) {
	workers := make([]iTopicWorker, 0, 16)
	defer func() {
		for _, w := range workers {
			w.close()
		}
	}()

	registerWorker := func(w iTopicWorker) {
		workers = append(workers, w)
	}

	view, err := newTopicWorker(
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
	registerWorker(&view)

	like, err := newTopicWorker(
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
	registerWorker(&like)

	comment, err := newTopicWorker(
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
	registerWorker(&comment)

	registration, err := newTopicWorker(
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
	registerWorker(&registration)

	posts, err := newTopicWorker(
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
	registerWorker(&posts)

	result := Workers{
		workers: workers,
	}
	workers = nil
	return result, nil
}

func (r *Workers) Start(ctx context.Context) {
	for _, w := range r.workers {
		w.start(ctx)
	}
}

func (r *Workers) Close() error {
	errs := make([]error, len(r.workers))
	for i, w := range r.workers {
		errs[i] = w.close()
	}

	return errors.Join(errs...)
}
