package main

import (
	"goschedule/http"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	server := http.NewServer()

	server.Addr = "localhost"
	server.Port = "4000"

	if err := server.Open(); err != nil {
		zap.S().Fatalf("%+v", err)
	}
}
