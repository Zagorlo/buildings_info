package main

import (
	"buildings_info/config"
	"buildings_info/logging"
	"buildings_info/server"
	"context"
	"runtime"
)

func init() {
	logging.InitLogger()
	logging.Logging.Info("Go version: %s\n", runtime.Version())
}

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	server.Launch(config.InitConfig(ctx), ctx, cancel)
}
