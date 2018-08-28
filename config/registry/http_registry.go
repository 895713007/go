package registry

import (
	"sync"
	"net/http"
	"time"
	"errors"
	"io/ioutil"
	"bytes"
	"fmt"
)

type httpRegistry struct {
	Host       string
	HttpClient *http.Client
	TTL        time.Duration
	sync.Mutex
}

func NewHttpRegistry(opts ...Option) Registry {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	//ttl minimum 5 seconds
	minTTL := time.Second * 5
	if options.TTL > minTTL {
		minTTL = options.TTL
	}

	timeout := time.Second * 3
	if options.Timeout > 0 {
		timeout = options.Timeout
	}

	return &httpRegistry{
		Host:       options.Host,
		TTL:        minTTL,
		HttpClient: &http.Client{Timeout: timeout},
	}
}

func (c *httpRegistry) Get(name string) (string, error) {
	b, err := c.request("GET", "/get", nil)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *httpRegistry) Set(name string, value string) error {
	return errors.New("todo")
}

func (c *httpRegistry) request(method string, path string, data []byte) ([]byte, error) {
	client := &http.Client{Timeout: time.Second * 5}

	url := c.Host + path
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
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
