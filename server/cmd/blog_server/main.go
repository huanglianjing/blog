package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/huanglianjing/blog/server/internal/common"
	"github.com/huanglianjing/blog/server/internal/model"
	"github.com/huanglianjing/blog/server/internal/router"
)

func main() {
	cfgPath := flag.String("c", "config/config.yaml", "配置文件路径")
	flag.Parse()

	cfg, err := common.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if err := model.InitDB(cfg.Database.Path); err != nil {
		log.Fatalf("init db: %v", err)
	}
	warnMissingContent()

	engine := router.New()
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}

// warnMissingContent 检查数据库中是否有搜索用的正文纯文本。
// 旧版 article_converter 写的库没有这一列的数据，此时正文搜索恒为空且不报错，
// 很难自查，故启动时给出提示。
func warnMissingContent() {
	total, err := model.CountArticles()
	if err != nil || total == 0 {
		return
	}
	withContent, err := model.CountArticlesWithContent()
	if err != nil {
		return
	}
	if withContent == 0 {
		log.Printf("警告: 共 %d 篇文章但均无正文纯文本，正文搜索将无结果，请重新运行 article_converter", total)
	}
}
