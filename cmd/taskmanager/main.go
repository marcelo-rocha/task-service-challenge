package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"

	conf "github.com/ardanlabs/conf/v2"
	"github.com/marcelo-rocha/task-service-challenge/server"
)

type Config struct {
	Server      server.ServerCfg
	Development bool `conf:"env:DEVELOPMENT,default:true"`
}

var cfg Config

func main() {
	_, err := conf.Parse("", &cfg)
	if err != nil {
		instructions, _ := conf.UsageInfo("", &cfg)
		fmt.Println(instructions)
		os.Exit(1)
	}

	var logger *zap.Logger
	if cfg.Development {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()

	srv := server.New(&cfg.Server, logger)
	if err := srv.Init(context.Background()); err != nil {
		logger.Fatal("initialization fail", zap.Error(err))
	}

	go func() {
		if err := srv.Run(); err != nil {
			logger.Info("run terminated", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logger.Info("Shutting down")
	os.Exit(0)
}
