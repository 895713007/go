package driver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mytokenio/go/log"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	CodeSuccess = 0
)

var ConfigServer = os.Getenv("CONFIG_SERVER")

type httpDriver struct {
	Host       string
	HttpClient *http.Client
	sync.Mutex
}

type Request struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Comment   string `json:"comment"`
	CreatedBy string `json:"created_by"`
}

type Config struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Comment   string `json:"comment"`
	UpdatedBy string `json:"updated_by"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
func (c Config) toMetadata() map[string]string {
	return map[string]string {
		"comment": c.Comment,
		"updated_by": c.UpdatedBy,
		"updated_at": c.UpdatedAt,
		"created_at": c.CreatedAt,
	}
}

type Response struct {
	Code      int     `json:"code"`
	Msg       string  `json:"msg"`
	Timestamp string  `json:"timestamp"`
	Data      *Config `json:"data"`
}

type ListResponse struct {
	Code      int       `json:"code"`
	Msg       string    `json:"msg"`
	Timestamp string    `json:"timestamp"`
	Data      []*Config `json:"data"`
}

func NewHttpDriver(opts ...Option) Driver {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	timeout := time.Second * 3
	if options.Timeout > 0 {
		timeout = options.Timeout
	}

	if options.Host == "" {
		options.Host = ConfigServer
	}

	return &httpDriver{
		Host:       options.Host,
		HttpClient: &http.Client{Timeout: timeout},
	}
}

func (d *httpDriver) List() ([]*Value, error) {
	var vals []*Value

	uri := "/v1/config/item"
	b, err := d.request("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	rsp := &ListResponse{}
	json.Unmarshal(b, rsp)
	if rsp.Code != CodeSuccess {
		return nil, fmt.Errorf("get error %s", rsp.Msg)
	}

	for _, c := range rsp.Data {
		v := &Value{
			K: c.Key,
			V: []byte(c.Value),
			//Format:    "toml",
			Metadata: c.toMetadata(),
		}
		vals = append(vals, v)
	}

	return vals, nil
}

func (d *httpDriver) Get(key string) (*Value, error) {
	uri := fmt.Sprintf("/v1/config/item/%s", key)
	b, err := d.request("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	rsp := &Response{}
	if rsp.Code != CodeSuccess {
		return nil, fmt.Errorf("get error %s", rsp.Msg)
	}
	err = json.Unmarshal(b, rsp)
	if err != nil {
		log.Errorf("json unmarshal error %s", err)
		return nil, fmt.Errorf("json unmarshal error %s", err)
	}

	c := rsp.Data
	v := &Value{
		K: c.Key,
		V: []byte(c.Value),
		//Format:    "toml",
		Metadata: c.toMetadata(),
	}
	return v, nil
}

func (d *httpDriver) Set(value *Value) error {
	uri := "/v1/config/item"
	req := Request{
		Key:       value.K,
		Value:     value.String(),
		Comment:   "",
		CreatedBy: "sdk",
	}
	reqBytes, _ := json.Marshal(req)

	var b []byte
	existValue, err := d.Get(value.K)
	if existValue != nil {
		b, err = d.request("PATCH", uri, reqBytes)
	} else {
		b, err = d.request("POST", uri, reqBytes)
	}
	if err != nil {
		return fmt.Errorf("post failed %s", err)
	}

	rsp := &Response{}
	json.Unmarshal(b, rsp)
	if rsp.Code != CodeSuccess {
		return fmt.Errorf("response error %s", rsp.Msg)
	}

	return nil
}

func (d *httpDriver) request(method string, path string, data []byte) ([]byte, error) {
	if d.Host == "" {
		return nil, errors.New("config server host empty")
	}

	url := d.Host + path
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	resp, err := d.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		err = nil
	} else {
		err = errors.New(fmt.Sprintf("http status code %d, body %s", resp.StatusCode, respBody))
	}
	return respBody, err
}

func (d *httpDriver) String() string {
	return "http"
}
