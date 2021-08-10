package tests

import (
	"buildings_info/models"
	"context"
	"fmt"
)

func testCfg(addr, port string, ctx context.Context) *models.Config {
	return &models.Config{
		Postgres: models.Postgres{
			Addr:            fmt.Sprintf("%s:%s", addr, port),
			Network:         "tcp",
			User:            "",
			Password:        "",
			Database:        "",
			ApplicationName: "",
			DialTimeout:     5000000000,
			ReadTimeout:     5000000000,
			WriteTimeout:    5000000000,
			MaxRetries:      3,
			MinRetryBackoff: 1000000000,
			MaxRetryBackoff: 5000000000,
			PoolSize:        5,
			MinIdleConns:    2,
			MaxConnAge:      60000000000,
			PoolTimeout:     5000000000,
			IdleTimeout:     5000000000,
		},
		BuildingsCache: models.BuildingsCache{
			10,
		},
		Context: ctx,
	}
}
