package buildings

import (
	"buildings_info/logging"
	"buildings_info/models"
	"context"
)

func (bm *BuildingsModule) UpdateBuildingInfo(requestCtx context.Context, building *models.Building) *logging.ErrorContainer {
	buildingChanges, ec := bm.driver.UpdateBuildingInfoP(requestCtx, building)
	if ec.NotNil() {
		return ec
	}

	bm.cache.PrependItems(requestCtx, buildingChanges)

	return logging.NilErrorContainerVar
}

func (bm *BuildingsModule) InsertBuildingInfo(requestCtx context.Context, building *models.Building) *logging.ErrorContainer {
	changes, ec := bm.driver.InsertBuildingInfoP(requestCtx, building)
	if ec.NotNil() {
		return ec
	}

	bm.cache.PrependItems(requestCtx, changes)

	return logging.NilErrorContainerVar
}

func (bm *BuildingsModule) GetBuildingsChangesInfo() []byte {
	return bm.cache.RetrieveItems()
}

func (bm *BuildingsModule) FillBuildingsCache(requestCtx context.Context) *logging.ErrorContainer {
	changes, ec := bm.driver.GetBuildingsChangesInfoP(requestCtx)
	if ec.NotNil() {
		return ec
	}

	bm.cache.PrependItems(requestCtx, changes)

	return logging.NilErrorContainerVar
}
