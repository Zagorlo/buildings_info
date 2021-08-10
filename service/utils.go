package service

import (
	"buildings_info/consts"
	"buildings_info/logging"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func refuse(w http.ResponseWriter, ec *logging.ErrorContainer) {
	http.Error(w, fmt.Sprintf("Refuse: %v, UUID: %s", ec.Err(), ec.UUID()), ec.Status())

	return
}

func answer(requestCtx context.Context, w http.ResponseWriter, ec *logging.ErrorContainer, value interface{}, method string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(func() int {
		if ec != nil && ec.Status() > 0 {
			return ec.Status()
		}

		return http.StatusOK
	}())

	if value != nil {
		err := json.NewEncoder(w).Encode(value)
		if err != nil {
			logging.NewErrorContainer(requestCtx, err, consts.JsonMarshalErrorStatus)
		}
	}

	logging.Logging.Info("Answer; method: %s; date: %s",
		method, time.Now().Format(time.RFC3339))
}

func answerUnencoded(requestCtx context.Context, w http.ResponseWriter, ec *logging.ErrorContainer, value []byte, method string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(func() int {
		if ec != nil && ec.Status() > 0 {
			return ec.Status()
		}

		return http.StatusOK
	}())

	_, err := w.Write(value)
	if err != nil {
		logging.NewErrorContainer(requestCtx, err, consts.WriteBodyStatus)
	}

	logging.Logging.Info("Answer; method: %s; date: %s",
		method, time.Now().Format(time.RFC3339))
}

func unmarshaller(requestCtx context.Context, r io.Reader, dest interface{}) *logging.ErrorContainer {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return logging.NewErrorContainer(requestCtx, err, http.StatusBadRequest)
	}

	err = json.Unmarshal(bytes, &dest)
	if err != nil {
		return logging.NewErrorContainer(requestCtx, err, http.StatusBadRequest, string(bytes))
	}

	return logging.NilErrorContainerVar
}
