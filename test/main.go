package main

import (
	"github.com/mytokenio/go_sdk/log"
	"github.com/mytokenio/go_sdk/config"
	"time"
	"github.com/mytokenio/go_sdk/config/driver"
	"github.com/mytokenio/go_sdk/registry"
	"google.golang.org/grpc/metadata"
	"encoding/json"
)

func main() {
	testService()
}

func testService() {
	s := registry.Service{
		Name:"test",
		Metadata: metadata.Pairs("kk", "vv", "aa", "bb", "aa", "cc"),
		Nodes: []registry.Node{
			{"test", "127.0.0.1", 12345},
		},
	}
	b, e := json.Marshal(s)
	if e != nil {
		log.Errorf("error %s", e)
	} else {
		log.Infof("service %s", b)
	}
}


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

func testConfig() {
	myConfigJson := `
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
	log.Infof("hello")

	mc := &MyConfig{}
	c := config.NewConfig(
		config.Service("mt.user"),
		config.TTL(time.Second * 10),
		config.Driver(driver.NewMockDriver()),
	)

	// config with http registry
	//reg := registry.NewHttpRegistry(registry.Host("http://xxx.com"))
	//c := config.NewConfig(
	//	config.Service("mt.user"),
	//	config.TTL(time.Second * 10),
	//	config.Registry(reg),
	//)

	c.Driver.Set(c.Service, []byte(myConfigJson))
	c.BindJSON(mc)

	b, err := c.GetServiceConfig()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Infof("get key %s %s", c.Service, b)
	log.Infof("MyConfig %v", mc)
}
