package config

import (
	"buildings_info/models"
	"context"
	"log"

	_ "github.com/lib/pq"
)

func InitConfig(ctx context.Context) *models.Config {
	cfg := getConfig(ctx)

	log.Println(cfg)

	return cfg
}
