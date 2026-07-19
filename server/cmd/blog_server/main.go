package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/huanglianjing/blog/server/internal/config"
	"github.com/huanglianjing/blog/server/internal/model"
	"github.com/huanglianjing/blog/server/internal/router"
)

func main() {
	cfgPath := flag.String("c", "config/config.yaml", "配置文件路径")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if err := model.InitDB(cfg.Database.Path); err != nil {
		log.Fatalf("init db: %v", err)
	}

	engine := router.New()
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
