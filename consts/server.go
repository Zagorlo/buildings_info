package consts

import "time"

const (
	ListenPortHTTP = ":1001"

	ReadTimeout  = 5 * time.Second
	WriteTimeout = 5 * time.Second
	IdleTimeout  = 60 * time.Second

	MaxHeaderBytes = 16 * 1024 * 1024

	ServerShutdownAwait = 0 * time.Second
)
