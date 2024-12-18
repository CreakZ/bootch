package main

import (
	"bootch/internal/cfg"
	"bootch/internal/routing"
	"bootch/pkg/cache"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	cfg := cfg.MustInitConfig()

	_ = cache.NewCache(cfg)

	routing.InitRouting(&r.RouterGroup)

	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
