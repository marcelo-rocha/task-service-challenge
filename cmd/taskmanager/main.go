package main

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"

	conf "github.com/ardanlabs/conf/v2"
	"github.com/marcelo-rocha/task-service-challenge/domain/task"
	"github.com/marcelo-rocha/task-service-challenge/persistence"
	"github.com/marcelo-rocha/task-service-challenge/server"
)

type Config struct {
	Server server.ServerCfg
	DBUrl  string
}

var cfg Config

func main() {
	_, err := conf.Parse("", &cfg)
	if err != nil {
		instructions, _ := conf.UsageInfo("", &cfg)
		fmt.Println(instructions)
		os.Exit(1)
	}
	logger, _ := zap.NewDevelopment()

	connection, err := persistence.NewConnection(context.Background(), cfg.DBUrl)
	if err != nil {
		fmt.Println("failed to connect to database", err)
		os.Exit(2)
	}
	defer connection.Close()

	tasksRepository := persistence.NewTasks(connection, logger)

	userCases := server.UserCases{
		NewTaskUseCase:      task.NewTaskUseCase{Persistence: tasksRepository},
		ListTasksUseCase:    task.ListTasksUseCase{Persistence: tasksRepository},
		FinalizeTaskUseCase: task.FinalizeTaskUseCase{Persistence: tasksRepository},
	}

	server.Run(&cfg.Server, logger, &userCases)
}
