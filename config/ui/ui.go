package main

import (
	"github.com/mytokenio/go/web"
	"github.com/mytokenio/go/log"
	"github.com/mytokenio/go/registry"
	"flag"
)

func main() {
	listen := flag.String("listen", ":5556", "server listen on ..")
	flag.Parse()

	//default mock registry, TODO
	reg := registry.NewRegistry()

	service := web.NewService(
		web.Name("config-ui"),
		web.Address(*listen),
		web.Handler(getHandler()),
		web.Registry(reg),
	)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
