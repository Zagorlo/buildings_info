package service

import (
	"buildings_info/logging"
	"buildings_info/models"
	"net/http"
)

func (s *Service) GetBuildingsChangesInfoSV(w http.ResponseWriter, r *http.Request) {
	answerUnencoded(r.Context(), w, logging.NilErrorContainerVar, s.GetBuildingsChangesInfo(), "GetBuildingsChangesInfoSV")
}

func (s *Service) UpdateBuildingInfoSV(w http.ResponseWriter, r *http.Request) {
	var building models.Building
	ec := unmarshaller(r.Context(), r.Body, &building)
	if ec.NotNil() {
		refuse(w, ec)
	}

	ec = s.UpdateBuildingInfo(r.Context(), &building)
	if ec.NotNil() {
		refuse(w, ec)
	}

	answer(r.Context(), w, logging.NilErrorContainerVar, building, "UpsertBuildingInfoSV")
}

func (s *Service) InsertBuildingInfoSV(w http.ResponseWriter, r *http.Request) {
	var building models.Building
	ec := unmarshaller(r.Context(), r.Body, &building)
	if ec.NotNil() {
		refuse(w, ec)
	}

	ec = s.InsertBuildingInfo(r.Context(), &building)
	if ec.NotNil() {
		refuse(w, ec)
	}

	answer(r.Context(), w, logging.NilErrorContainerVar, building, "InsertBuildingInfoSV")
}
