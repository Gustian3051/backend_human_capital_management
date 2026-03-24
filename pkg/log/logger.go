package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init(isDebug bool) {
	var err error

	if isDebug {
		Log, err = zap.NewDevelopment()
	} else {
		Log, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}