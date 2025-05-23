package bootstrap

import (
	"context"
	"gateway/app/repository"
)

func (bs *Bootstrap) initRepository() {
	bs.Repository = repository.New(bs.Database)

	ctx := context.Background()
	bs.Repository.Example.Migrate(ctx)
	bs.Database.BuildIndexes(ctx, "app/model")
}
