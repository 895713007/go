package log

import (
	"os"
	"github.com/sirupsen/logrus"
	"github.com/rs/xid"
	"github.com/json-iterator/go"
)

const (
	PanicLevel logrus.Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

const typeField = "log_type"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var log = logrus.New()
var extra map[string]interface{}
var uniqueId string

func init() {
	log.SetOutput(os.Stdout)

	if _level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		logrus.Infof("set level %s %s", _level, os.Getenv("LOG_LEVEL"))
		log.SetLevel(_level)
	} else {
		log.SetLevel(DebugLevel)
	}
	Debugf("level %s", log.Level)

	server := os.Getenv("LOG_SERVER")
	if server != "" {
		log.AddHook(NewEsLogHook(server))
	}

	SetExtra(map[string]interface{}{
		"command": os.Args[0],
	})
}

func SetLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetExtra(h map[string]interface{}) {
	extra = h
	RefreshUniqueId()
}

func RefreshUniqueId() {
	uniqueId = xid.New().String()
}

func Type(typ string) *logrus.Entry {
	return log.WithField(typeField, typ)
}

func WithField(key string, value interface{}, typ string) *logrus.Entry {
	return log.WithField(typeField, typ).WithField(key, value)
}

func WithFields(fields logrus.Fields, typ string) *logrus.Entry {
	fields[typeField] = typ
	return log.WithFields(fields)
}

func Info(v ...interface{}) {
	log.Info(v...)
}

func Infof(format string, v ...interface{}) {
	log.Infof(format, v...)
}

func Debug(v ...interface{}) {
	log.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func Warn(v ...interface{}) {
	log.Warning(v...)
}

func Warnf(format string, v ...interface{}) {
	log.Warningf(format, v...)
}

func Warning(v ...interface{}) {
	log.Warning(v...)
}

func Warningf(format string, v ...interface{}) {
	log.Warningf(format, v...)
}

func Error(v ...interface{}) {
	log.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}
