package driver

import (
	"sync"
	"net/http"
	"time"
	"errors"
	"io/ioutil"
	"bytes"
	"fmt"
	"encoding/json"
	"os"
)

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

type Response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Data struct {
		Key       string `json:"key"`
		Value     string `json:"value"`
		Comment   string `json:"comment"`
		UpdatedBy string `json:"updated_by"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
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
		options.Host = os.Getenv("CONFIG_SERVER")
	}

	return &httpDriver{
		Host:       options.Host,
		HttpClient: &http.Client{Timeout: timeout},
	}
}

func (d *httpDriver) List() ([]*Value, error) {
	var vals []*Value
	return vals, nil
}

func (c *httpDriver) Get(key string) (*Value, error) {
	uri := fmt.Sprintf("/v1/config/item/%s", key)
	b, err := c.request("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	rsp := &Response{}
	json.Unmarshal(b, rsp)
	if rsp.Code != 0 {
		return nil, fmt.Errorf("get error %s", rsp.Message)
	}

	v := &Value{
		K: key,
		V: b,
		Timestamp: rsp.Data.UpdatedAt.Unix(),
		Format: "json",
		Metadata: map[string]string{},
	}
	return v, nil
}

func (c *httpDriver) Set(value *Value) error {
	uri := "/v1/config/item"
	req := Request{
		Key: value.K,
		Value: value.String(),
		Comment: "",
		CreatedBy: "sdk",
	}
	reqBytes, _ := json.Marshal(req)
	b, err := c.request("POST", uri, reqBytes)
	if err != nil {
		return fmt.Errorf("post failed %s", err)
	}

	rsp := &Response{}
	json.Unmarshal(b, rsp)
	if rsp.Code != 0 {
		return fmt.Errorf("response error %s", rsp.Message)
	}

	return nil
}

func (c *httpDriver) request(method string, path string, data []byte) ([]byte, error) {
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

func (c *httpDriver) String() string {
	return "http"
}