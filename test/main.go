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
	testConfig()
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
//	myConfigJson := `
//{
//	"api": "http://api.mytokenapi.com",
//	"db": {
//		"host": "localhost",
//		"user": "root",
//		"password": "",
//		"name": "mytoken"
//	},
//	"log_servers": ["127.0.0.1:12333", "127.0.0.1:12334"]
//}
//`
	log.Infof("hello")

	d := driver.NewFileDriver(driver.Path("./config.json"))

	mc := &MyConfig{}
	c := config.NewConfig(
		config.Service("mt.user"),
		config.TTL(time.Second * 10),
		config.Driver(d),
		config.OnChanged(func(c *config.Config) error {
			err := c.BindJSON(mc)
			if err != nil {
				return err
			}

			//panic("test panic")
			log.Infof("service config changed %v", mc)
			//TODO

			return nil
		}),
	)

	//// config with http driver
	//c := config.NewConfig(
	//	config.Service("mt.user"),
	//	config.TTL(time.Second * 10), //cache ttl
	//	config.Driver(
	//		driver.NewHttpDriver(
	//			driver.Host("http://xxx.com"),
	//			driver.Timeout(time.Second * 3),
	//		),
	//	),
	//)

	//c.Driver.Set(c.Service, []byte(myConfigJson))
	c.BindJSON(mc)

	b, err := c.GetServiceConfig()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Infof("get key %s %s", c.Service, b)
	log.Infof("MyConfig %v", mc)

	time.Sleep(time.Minute)
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
