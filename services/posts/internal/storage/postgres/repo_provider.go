package postgres

import (
	"context"
	"soa-socialnetwork/services/posts/internal/repo"
)

type repoProvider struct {
	ctx   context.Context
	scope pgxScope
}

func (p *repoProvider) Pages() repo.PagesRepository {
	return pagesRepo{
		ctx:   p.ctx,
		scope: p.scope,
	}
}
func (p *repoProvider) Posts() repo.PostsRepository {
	return postsRepo{
		ctx:   p.ctx,
		scope: p.scope,
	}
}
func (p *repoProvider) Comments() repo.CommentsRepository {
	return commentsRepo{
		ctx:   p.ctx,
		scope: p.scope,
	}
}
func (p *repoProvider) Metrics() repo.MetricsRepository {
	return metricsRepo{
		ctx:   p.ctx,
		scope: p.scope,
	}
}
func (p *repoProvider) Outbox() repo.OutboxRepository {
	return outboxRepo{
		ctx:   p.ctx,
		scope: p.scope,
	}
}
