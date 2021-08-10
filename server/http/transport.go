package httpServer

import (
	"buildings_info/consts"
	"buildings_info/service"
	"context"
	uuid "github.com/satori/go.uuid"
	"net/http"

	"github.com/gorilla/mux"
)

func middleware() func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), consts.ContextUUIDKey, uuid.NewV4())
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func NewServer(serverErrors chan error, ser *service.Service) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/health", healthHandler).Methods(http.MethodGet)

	//api := router.PathPrefix("/api")
	apiV1 := router.PathPrefix("/api/v1").Subrouter()
	apiV1.Use(middleware())
	apiV1.Path("/building/update").Methods(http.MethodPut).HandlerFunc(ser.UpdateBuildingInfoSV)
	apiV1.Path("/building/insert").Methods(http.MethodPost).HandlerFunc(ser.InsertBuildingInfoSV)
	apiV1.Path("/changes").Methods(http.MethodGet).HandlerFunc(ser.GetBuildingsChangesInfoSV)

	server := &http.Server{
		Addr:           consts.ListenPortHTTP,
		Handler:        serveOrigin([]string{"*"}, router),
		MaxHeaderBytes: consts.MaxHeaderBytes,
		ReadTimeout:    consts.ReadTimeout,
		WriteTimeout:   consts.WriteTimeout,
		IdleTimeout:    consts.IdleTimeout,
	}

	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	return server
}
