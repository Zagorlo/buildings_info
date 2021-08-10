package models

import (
	"context"
	"time"
)

type Config struct {
	Postgres       Postgres       `toml:"Postgres"`
	BuildingsCache BuildingsCache `toml:"BuildingsCache"`
	Context        context.Context
}

type Postgres struct {
	Addr            string        `toml:"Addr"`
	Network         string        `toml:"Network"`
	User            string        `toml:"User"`
	Password        string        `toml:"Password"`
	Database        string        `toml:"Database"`
	ApplicationName string        `toml:"ApplicationName"`
	DialTimeout     time.Duration `toml:"DialTimeout"`
	ReadTimeout     time.Duration `toml:"ReadTimeout"`
	WriteTimeout    time.Duration `toml:"WriteTimeout"`
	MaxRetries      int           `toml:"MaxRetries"`
	MinRetryBackoff time.Duration `toml:"MinRetryBackoff"`
	MaxRetryBackoff time.Duration `toml:"MaxRetryBackoff"`
	PoolSize        int           `toml:"PoolSize"`
	MinIdleConns    int           `toml:"MinIdleConns"`
	MaxConnAge      time.Duration `toml:"MaxConnAge"`
	PoolTimeout     time.Duration `toml:"PoolTimeout"`
	IdleTimeout     time.Duration `toml:"IdleTimeout"`
	RefreshAwait    time.Duration `toml:"RefreshAwait"`
	RefreshCheck    time.Duration `toml:"RefreshCheck"`
}

type BuildingsCache struct {
	Size int64 `toml:"Size"`
}
