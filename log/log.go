package log

import (
	"fmt"
	"time"
)

var (
	c  *Config
	hs = make([]handle, 0)
)

type handle interface {
	info(format string, args ...interface{})
	warn(format string, args ...interface{})
	debug(format string, args ...interface{})
	error(format string, args ...interface{})
}

type xtime time.Time

// Config log config.
type Config struct {
	Stdout bool
	Dir    string
	Agent  *AgentConfig
}

// AgentConfig agent config.
type AgentConfig struct {
	TaskID  string
	Proto   string
	Addr    string
	Buffer  int
	Chan    int
	Timeout time.Duration
}

type D struct {
	Index    string   `json:"@index"`
	Type     string   `json:"@type"`
	Time     xtime    `json:"datetime"`
	UniqueID string   `json:"unique_id"`
	UID      int      `json:"uid"`
	Info     *logInfo `json:"info"`
}

type logInfo struct {
	Host    string `json:"host"`
	Extra   string `json:"extra"`
	Message string `json:"message"`
	Context string `json:"context"`
}

func (xt xtime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(xt).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

// Init create logger with context.
func Init(conf *Config) {
	if conf.Stdout == true {
		hs = append(hs, stdInit())
	}
	if isDir(conf.Dir) {
		hs = append(hs, fileInit(conf))
	}
}

func Info(format string, args ...interface{}) {
	for _, h := range hs {
		h.info(format, args...)
	}
}

func Warn(format string, args ...interface{}) {
	for _, h := range hs {
		h.warn(format, args...)
	}
}
func Debug(format string, args ...interface{}) {
	for _, h := range hs {
		h.debug(format, args...)
	}
}

func Error(format string, args ...interface{}) {
	for _, h := range hs {
		h.error(format, args...)
	}
}
