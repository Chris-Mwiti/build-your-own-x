package logger

import (
	"log"

	"go.uber.org/zap"
)

func Init() (*zap.Logger, *zap.SugaredLogger){
	
	devlogger,err := zap.NewDevelopment()
	defer devlogger.Sync()
	if err != nil {
		log.Fatalf("fatal: could not setup logger: %s", err.Error())
	}

	sugarLogger := devlogger.Sugar()

	return devlogger, sugarLogger
}

