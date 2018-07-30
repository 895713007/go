package log

import (
	"fmt"
	"runtime"

	pkgerr "github.com/pkg/errors"
)

const (
	_timeFormat = "2006-01-02T15:04:05.999999"

	_levelValue = "level_value"
	_level      = "level"
	_time       = "time"
	_title      = "title"
	_source     = "source"
	_log        = "log"
	_appID      = "app_id"
	_instanceID = "instance_id"
	_tid        = "traceid"
	_ts         = "ts"
	_caller     = "caller"
)

type Handler interface {
	Log(lv Level, f D)
	Close() error
}

// Handlers .
type Handlers []Handler

// Log handlers logging.
func (hs Handlers) Log(lv Level, f D) {
	for _, h := range hs {
		h.Log(lv, f)
	}
}

// Close close resource.
func (hs Handlers) Close() (err error) {
	for _, h := range hs {
		if e := h.Close(); e != nil {
			err = pkgerr.WithStack(e)
		}
	}
	return
}

// funcName get func name.
func funcName() (name string) {
	if pc, _, lineNo, ok := runtime.Caller(3); ok {
		name = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineNo)
	}
	return
}
