package postgres

import (
	"context"
	"net/http"

	"buildings_info/consts"
	"buildings_info/logging"
	"buildings_info/models"

	"github.com/go-pg/pg/v10"
	"github.com/jinzhu/copier"
)

func NewBuildingsPostgresDriver(ctx context.Context, pgConn *pg.DB) *BuildingsPostgresDriver {
	return &BuildingsPostgresDriver{ctx: ctx, pgConn: pgConn}
}

type BuildingsPostgresDriver struct {
	ctx    context.Context
	pgConn *pg.DB
}

func (bpd *BuildingsPostgresDriver) UpdateBuildingInfoP(requestCtx context.Context, building *models.Building) ([]models.CacheItem, *logging.ErrorContainer) {
	if building == nil {
		return nil, logging.NewErrorContainer(requestCtx, consts.ErrorNilPointerReference, http.StatusBadRequest)
	}

	var newBuilding = new(models.Building)
	copier.Copy(newBuilding, building)

	_, err := bpd.pgConn.Model(building).WherePK().Returning(
		`
			(select id from buildings where id = ?id),
			(select name from buildings where id = ?id),
			(select floors_count from buildings where id = ?id),
			(select parking_count from buildings where id = ?id),
			(select parking_available from buildings where id = ?id)`,
	).Update()
	if err != nil {
		return nil, logging.NewErrorContainer(requestCtx, err, http.StatusInternalServerError, building)
	}

	changes := newBuilding.RetrieveChanges(building)

	return changes, logging.NilErrorContainerVar
}

func (bpd *BuildingsPostgresDriver) InsertBuildingInfoP(requestCtx context.Context, building *models.Building) ([]models.CacheItem, *logging.ErrorContainer) {
	if building == nil {
		return nil, logging.NewErrorContainer(requestCtx, consts.ErrorNilPointerReference, http.StatusBadRequest)
	}

	building.PrepareValues()

	_, err := bpd.pgConn.Model(building).Returning("id").Insert()
	if err != nil {
		return nil, logging.NewErrorContainer(requestCtx, err, http.StatusInternalServerError, building)
	}

	return building.RetrieveChanges(nil), logging.NilErrorContainerVar
}

func (bpd *BuildingsPostgresDriver) GetBuildingsChangesInfoP(requestCtx context.Context) ([]models.CacheItem, *logging.ErrorContainer) {
	var buildings []*models.BuildingUpdate
	var changes []models.CacheItem
	err := bpd.pgConn.Model(&buildings).Limit(10).Order("removed_at desc").Select()
	if err != nil {
		return nil, logging.NewErrorContainer(requestCtx, err, http.StatusNoContent)
	}

	for i := range buildings {
		changes = append(changes, buildings[i].RetrieveBuildingUpdates()...)
	}

	return changes, logging.NilErrorContainerVar
}
