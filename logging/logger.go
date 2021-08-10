package logging

import (
	"buildings_info/consts"
	"context"
	"encoding/json"
	uuid2 "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"time"
)

var Logging *Logger

type Logger struct {
	logger *zap.SugaredLogger
}

func InitLogger() {
	Logging = new(Logger)
	Logging.logger = new(zap.SugaredLogger)
	Logging.logger = serviceLogger()
}

func serviceLogger() *zap.SugaredLogger {
	loggerCfg := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(loggerCfg, &cfg); err != nil {
		panic(err)
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Infof(msg, args...)
	} else {
		l.logger.Infow(msg)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Errorf(msg, args...)
	} else {
		l.logger.Errorw(msg)
	}
}

func (l *Logger) Sync() {
	_ = l.logger.Sync()
}

func GetRequestUUID(requestCtx context.Context) uuid2.UUID {
	uuidCtx := requestCtx.Value(consts.ContextUUIDKey)
	uuidParsed, parsed := uuidCtx.(uuid2.UUID)
	if !parsed {
		Logging.Error("unable to parse uuid at: %s", time.Now().Format(time.RFC3339))
	}

	return uuidParsed
}
