package registry

import (
	"sync"
	"net/http"
	"time"
	"errors"
	"io/ioutil"
	"bytes"
	"fmt"
	"encoding/json"
)

type httpRegistry struct {
	Host       string
	HttpClient *http.Client
	sync.Mutex
}

type Response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Data struct {
		Key       string `json:"key"`
		Value     string `json:"value"`
		Comment   string `json:"comment"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
}

func NewHttpRegistry(opts ...Option) Registry {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	timeout := time.Second * 3
	if options.Timeout > 0 {
		timeout = options.Timeout
	}

	return &httpRegistry{
		Host:       options.Host,
		HttpClient: &http.Client{Timeout: timeout},
	}
}

func (c *httpRegistry) Get(key string) ([]byte, error) {
	uri := fmt.Sprintf("/v1/config/item/%s", key)
	b, err := c.request("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	rsp := &Response{}
	json.Unmarshal(b, rsp)
	return []byte(rsp.Data.Value), nil
}

func (c *httpRegistry) Set(name string, value []byte) error {
	return errors.New("todo")
}

func (c *httpRegistry) request(method string, path string, data []byte) ([]byte, error) {
	if c.Host == "" {
		return nil, errors.New("config server host empty")
	}

	url := c.Host + path
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		err = nil
	} else {
		err = errors.New(fmt.Sprintf("http get errcode %d, errmsg %s", resp.StatusCode, respBody))
	}
	return respBody, err
}

func (c *httpRegistry) String() string {
	return "http"
}