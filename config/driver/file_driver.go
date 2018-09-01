package driver

import (
	"os"
	"io/ioutil"
	"errors"
)

type fileDriver struct {
	path string
}

func NewFileDriver(opts ...Option) Driver {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	return &fileDriver{
		path: options.Path,
	}
}

func (d *fileDriver) Get(key string) ([]byte, error) {
	path := key
	if d.path != "" {
		path = d.path
	}

	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	b, err := ioutil.ReadAll(fh)
	if err != nil {
		return nil, err
	}

	//info, err := fh.Stat()
	//if err != nil {
	//	return nil, err
	//}
	return b, nil
}

func (d *fileDriver) Set(key string, value []byte) error {
	return errors.New("TODO")
}

func (d *fileDriver) String() string {
	return "file"
}