package main

import (
	"fmt"
	"os"

	conf "github.com/ardanlabs/conf/v2"
	"github.com/marcelo-rocha/task-service-challenge/server"
)

type Config struct {
	Server server.ServerCfg
}

var cfg Config

func main() {
	_, err := conf.Parse("", &cfg)
	if err != nil {
		instructions, _ := conf.UsageInfo("", &cfg)
		fmt.Println(instructions)
		os.Exit(1)
	}
	server.Run(&cfg.Server)
}
