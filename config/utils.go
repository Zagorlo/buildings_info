package config

import (
	"buildings_info/models"
	"context"
	"github.com/BurntSushi/toml"
)

func getConfig(ctx context.Context) *models.Config {
	cfg := &models.Config{Context: ctx}
	path := "config/config.toml"

	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
