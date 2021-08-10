package buildings

import (
	"buildings_info/storages/buildings/storage/postgres"
	"buildings_info/storages/models/serviceCache"
	"context"
	"github.com/go-pg/pg/v10"
)

type BuildingsModule struct {
	ctx    context.Context
	driver *postgres.BuildingsPostgresDriver
	cache  *serviceCache.JsonOrderedCache
}

func (bm *BuildingsModule) Init(
	ctx context.Context,
	pgConn *pg.DB,
	suppliersCache *serviceCache.JsonOrderedCache,
) *BuildingsModule {
	bm.ctx = ctx
	bm.driver = postgres.NewBuildingsPostgresDriver(ctx, pgConn)
	bm.cache = suppliersCache

	return bm
}
