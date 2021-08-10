package server

import (
	"buildings_info/models"
	"buildings_info/storages/buildings"
	"buildings_info/storages/models/postgresDriver"
	"buildings_info/storages/models/serviceCache"
	"context"
	uuid2 "github.com/satori/go.uuid"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"buildings_info/consts"
	"buildings_info/logging"
	httpServer "buildings_info/server/http"
	"buildings_info/service"
)

func Launch(cfg *models.Config, ctx context.Context, ctxCancel func()) {
	ctx = context.WithValue(ctx, consts.ContextUUIDKey, uuid2.NewV4())

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	serverErrors := make(chan error, 1)
	driver := postgresDriver.NewPostgresDriver(&cfg.Postgres, ctx)
	ec := driver.PostgresOpen()
	if ec.NotNil() {
		panic(ec.Err())
	}

	ec = driver.PostgresCreations("./init.sql")
	if ec.NotNil() {
		panic(ec.Err())
	}

	buildingsModule := &buildings.BuildingsModule{}

	buildingsModule.Init(ctx, driver.GetConn(), serviceCache.NewJsonOrderedCache(cfg.BuildingsCache))

	ser := &service.Service{
		buildingsModule,
	}

	ec = ser.FillBuildingsCache(ctx)
	if ec.NotNil() {
		panic(ec.Err())
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				done := int64(1)
				go func() {
					defer atomic.AddInt64(&done, -1)
					driver.PostgresRefresh()
				}()

				for done > 0 {
					time.Sleep(250 * time.Millisecond)
				}
			}
		}
	}()

	server := httpServer.NewServer(serverErrors, ser)

	logging.Logging.Info("Сервер (http) стартанул по адресу: %v", server.Addr)

	select {
	case err := <-serverErrors:
		logging.Logging.Error("Ошибка сервера: %v", err)

	case <-osSignals:
		if err := server.Shutdown(ctx); err != nil {
			logging.Logging.Error("Вменяемого завершения не удалось: %v", err)
		}

		logging.Logging.Info("Завернули сервер")
	}

	beforeExit(ctxCancel, driver)
}

func beforeExit(ctxCancel func(), driver *postgresDriver.PostgresDriver) {
	ctxCancel()
	time.Sleep(consts.ServerShutdownAwait)

	driver.PostgresClose()

	logging.Logging.Info("Вменяемо вышли")
	logging.Logging.Sync()
}
