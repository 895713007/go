package main

import (
	"github.com/mytokenio/go_sdk/web"
	"github.com/mytokenio/go_sdk/log"
	"github.com/mytokenio/go_sdk/registry"
)

func main() {
	//default mock registry, TODO
	reg := registry.NewRegistry()

	service := web.NewService(
		web.Name("config-ui"),
		web.Address(":5556"),
		web.Handler(getHandler()),
		web.Registry(reg),
	)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
