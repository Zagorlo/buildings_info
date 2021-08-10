package postgresDriver

import (
	"buildings_info/logging"
	"buildings_info/models"
	"context"
	"github.com/go-pg/pg/v10"
	"io/ioutil"
	"net/http"
	"time"
)

func NewPostgresDriver(postgresConfig *models.Postgres, ctx context.Context) *PostgresDriver {
	return &PostgresDriver{
		postgresConfig: postgresConfig,
		ctx:            ctx,
	}
}

type PostgresDriver struct {
	pgConn *pg.DB

	postgresConfig *models.Postgres
	ctx            context.Context
}

func (pd *PostgresDriver) GetConn() *pg.DB {
	return pd.pgConn
}

func (pd *PostgresDriver) PostgresOpen() *logging.ErrorContainer {
	pgConn := pg.Connect(pd.pgOptions())

	pd.pgConn = pgConn

	err := pd.pgConn.Ping(pd.ctx)
	if err != nil {
		return logging.NewErrorContainer(pd.ctx, err, http.StatusInternalServerError, pd.postgresConfig)
	}

	return logging.NilErrorContainerVar
}

func (pd *PostgresDriver) PostgresCreations(filePath string) *logging.ErrorContainer {
	query, err := ioutil.ReadFile(filePath)
	if err != nil {
		return logging.NewErrorContainer(pd.ctx, err, http.StatusInternalServerError, pd.postgresConfig)
	}

	if _, err := pd.pgConn.Exec(string(query)); err != nil {
		return logging.NewErrorContainer(pd.ctx, err, http.StatusInternalServerError, pd.postgresConfig)
	}

	return logging.NilErrorContainerVar
}

func (pd *PostgresDriver) PostgresClose() *logging.ErrorContainer {
	err := pd.pgConn.Close()
	if err != nil {
		return logging.NewErrorContainer(pd.ctx, err, http.StatusExpectationFailed)
	}

	return logging.NilErrorContainerVar
}

func (pd *PostgresDriver) PostgresRefresh() {
	for {
		select {
		case <-pd.ctx.Done():
			return
		default:
		}

		if err := pd.pgConn.Ping(pd.ctx); err != nil {
			time.Sleep(pd.postgresConfig.RefreshAwait)
			logging.NewErrorContainer(pd.ctx, err, http.StatusInternalServerError, pd.postgresConfig.Addr)

			pd.pgConn = pg.Connect(pd.pgOptions())
		}

		time.Sleep(pd.postgresConfig.RefreshCheck)
	}
}

func (pd *PostgresDriver) pgOptions() *pg.Options {
	return &pg.Options{
		Network:         pd.postgresConfig.Network,
		Addr:            pd.postgresConfig.Addr,
		User:            pd.postgresConfig.User,
		Password:        pd.postgresConfig.Password,
		Database:        pd.postgresConfig.Database,
		ApplicationName: pd.postgresConfig.ApplicationName,
		DialTimeout:     pd.postgresConfig.DialTimeout,
		ReadTimeout:     pd.postgresConfig.ReadTimeout,
		WriteTimeout:    pd.postgresConfig.WriteTimeout,
		MaxRetries:      pd.postgresConfig.MaxRetries,
		MinRetryBackoff: pd.postgresConfig.MinRetryBackoff,
		MaxRetryBackoff: pd.postgresConfig.MaxRetryBackoff,
		PoolSize:        pd.postgresConfig.PoolSize,
		MinIdleConns:    pd.postgresConfig.MinIdleConns,
		MaxConnAge:      pd.postgresConfig.MaxConnAge,
		PoolTimeout:     pd.postgresConfig.PoolTimeout,
		IdleTimeout:     pd.postgresConfig.IdleTimeout,
	}
}
