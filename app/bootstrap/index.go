package bootstrap

import (
	"gateway/app/repository"
	"gateway/package/redis"
	"gateway/package/redsync"

	"github.com/RevenueMonster/sqlike/sqlike"
)

type Bootstrap struct {
	Database   *sqlike.Database
	Redsync    *redsync.Redsync
	Redis      *redis.Client
	Repository *repository.Repository
}

// New :
func New() *Bootstrap {
	bs := new(Bootstrap)
	bs.initMySQL()
	bs.initRedis()
	bs.initRepository()
	return bs
}
