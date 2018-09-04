package main

import (
	"github.com/mytokenio/go_sdk/log"
	"github.com/mytokenio/go_sdk/config"
	"time"
	"github.com/mytokenio/go_sdk/registry"
	"google.golang.org/grpc/metadata"
	"encoding/json"
	"github.com/mytokenio/go_sdk/config/driver"
)

func main() {
	testFileConfig()
	time.Sleep(time.Minute)
}

type MyConfig struct {
	API string `toml:"api"`
	DB struct {
		Host     string `toml:"host"`
		User     string `toml:"user"`
		Password string `toml:"password"`
		Name     string `toml:"name"`
	} `toml:"db"`
	LogServers []string `toml:"log_servers"`
}

func testFileConfig() {
	mc := &MyConfig{}

	config.NewConfig().Watch(func(c *config.Config) error {
		err := c.BindTOML(mc)
		if err != nil {
			log.Errorf("config bind error %s", err)
			return err
		}

		log.Infof("service config changed %v", mc)
		return nil
	})
	//
	//c := config.NewConfig(
	//	config.TTL(10 * time.Second),
	//	config.Driver(
	//		driver.NewFileDriver(driver.Path("config.toml")),
	//	),
	//)
	//c.Watch(func(c *config.Config) error {
	//	err := c.BindTOML(mc)
	//	if err != nil {
	//		log.Errorf("config bind error %s", err)
	//		return err
	//	}
	//
	//	log.Infof("service config changed %v", mc)
	//	return nil
	//})

	log.Infof("MyConfig %v", mc)
}

func testHttpConfig() {
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
	mc := &MyConfig{}
	c := config.NewConfig(
		config.Service("user"),
		config.TTL(time.Second*10), //cache ttl
		config.Driver(
			driver.NewHttpDriver(
				driver.Host("http://xxx.com"),
				driver.Timeout(time.Second*3),
			),
		),
	)

	value := driver.NewValue("mt.service."+c.Service, []byte(myConfigJson))
	c.Driver.Set(value)
	c.BindTOML(mc)
}

func testService() {
	s := registry.Service{
		Name:     "test",
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
