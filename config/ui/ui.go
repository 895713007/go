package main

import (
	"github.com/mytokenio/go_sdk/web"
	"github.com/mytokenio/go_sdk/log"
	"github.com/mytokenio/go_sdk/registry"
	"github.com/gin-gonic/gin"
)

func main() {
	//gin handler
	handler := gin.Default()
	servePages(handler)

	//default mock registry, TODO
	reg := registry.NewRegistry()

	service := web.NewService(
		web.Name("config-ui"),
		web.Address(":5556"),
		web.Handler(handler),
		web.Registry(reg),
	)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
