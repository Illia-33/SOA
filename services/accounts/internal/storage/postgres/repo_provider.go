package postgres

import (
	"context"
	"soa-socialnetwork/services/accounts/internal/repo"
)

type repoProvider struct {
	ctx   context.Context
	scope pgxScope
}

func (p *repoProvider) Accounts() repo.AccountsRepo {
	return accountsRepo{
		ctx:   p.ctx,
		scope: p.scope,
	}
}

func (p *repoProvider) Profiles() repo.ProfilesRepo {
	return profilesRepo{
		ctx:   p.ctx,
		scope: p.scope,
	}
}

func (p *repoProvider) ApiTokens() repo.ApiTokensRepo {
	return apiTokensRepo{
		ctx:   p.ctx,
		scope: p.scope,
	}
}

func (p *repoProvider) Outbox() repo.OutboxRepo {
	return outboxRepo{
		ctx:   p.ctx,
		scope: p.scope,
	}
}
