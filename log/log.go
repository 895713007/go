package log

import (
	"fmt"
	"os"
	"time"
)

// Config log config.
type Config struct {
	Stdout bool
	Host   string
	Dir    string
	Agent  *AgentConfig
	// VLevel Enable V-leveled logging at the specified _level.
	VLevel int32
	// Module=""
	// The syntax of the argument is a map of pattern=N,
	// where pattern is a literal file name (minus the ".go" suffix) or
	// "glob" pattern and N is a V _level. For instance:
	// [module]
	//   "service" = 1
	//   "dao*" = 2
	// sets the V _level to 2 in all Go files whose names begin "dao".
	Module map[string]int32
}

// D represents a map of entry _level data used for structured logging.
type D map[string]interface{}

var (
	h Handler
	c *Config
)

// Init create logger with context.
func Init(conf *Config) {
	var (
		hs Handlers
	)
	conf.Host = os.Getenv("CONF_HOSTNAME")
	if len(conf.Host) == 0 {
		hn, _ := os.Hostname()
		conf.Host = hn
	}
	c = conf
	if conf.Stdout {
		hs = append(hs, NewStdout())
	}
	if conf.Dir != "" {
		hs = append(hs, NewFile(conf.Dir))
	}
	if conf.Agent != nil {
		hs = append(hs, NewAgent(conf.Agent))
	}
	h = hs
}

// Info logs a message at the info log _level.
func Info(format string, args ...interface{}) {
	logf(infoLevel, format, args...)
}

// Warn logs a message at the warning log _level.
func Warn(format string, args ...interface{}) {
	logf(warnLevel, format, args...)
}

// Error logs a message at the error log _level.
func Error(format string, args ...interface{}) {
	logf(errorLevel, format, args...)
}

func logf(lv Level, format string, args ...interface{}) {
	if h == nil {
		return
	}
	now := time.Now()
	f := D{}
	f[_appID] = "app_ID"
	f[_instanceID] = c.Host
	f[_levelValue] = lv
	f[_level] = lv.String()
	f[_time] = now.Format(_timeFormat)
	f[_source] = funcName()
	f[_log] = fmt.Sprintf(format, args...)
	h.Log(lv, f)
}

// Close close resource.
func Close() (err error) {
	if h == nil {
		return
	}
	return h.Close()
}
