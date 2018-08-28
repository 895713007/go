package main

import (
	"github.com/mytokenio/go_sdk/log"
	"github.com/mytokenio/go_sdk/config"
	"time"
	"github.com/mytokenio/go_sdk/config/registry"
)

const MyConfigJson = `
{
	"api": "http://api.mytokenapi.com",
	"db": {
		"host": "localhost",
		"user": "root",
		"password": "",
		"name": "mytoken"
	},
	"log_servers": ["127.0.0.1:12333", "127.0.0.1:12334"]
}
`

type MyConfig struct {
	API string `json:"api"`
	DB struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"db"`
	LogServers []string `json:"log_servers"`
}

func main() {
	log.Infof("hello")

	mc := &MyConfig{}
	c := config.NewConfig(
		config.Service("mt.user"),
		config.TTL(time.Second * 10),
		config.Registry(registry.NewMockRegistry()),
	)

	// config with http registry
	//reg := registry.NewHttpRegistry(registry.Host("http://xxx.com"))
	//c := config.NewConfig(
	//	config.Service("mt.user"),
	//	config.TTL(time.Second * 10),
	//	config.Registry(reg),
	//)

	c.Registry.Set(c.Service, []byte(MyConfigJson))
	c.BindJSON(mc)

	b, err := c.Get()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Infof("get key %s %s", c.Service, b)
	log.Infof("MyConfig %v", mc)
}

func httpRegistry() {

	mc := &MyConfig{}


}