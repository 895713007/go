package log

import "github.com/Sirupsen/logrus"

type stdHandler struct {
	logger *logrus.Logger
}

func stdInit() *stdHandler {
	return &stdHandler{
		logger: new(logrus.Logger),
	}
}

func (fh *stdHandler) info(format string, args ...interface{}) {
	fh.logger.Infof(format, args...)
}

func (fh *stdHandler) warn(format string, args ...interface{}) {
	fh.logger.Warnf(format, args...)
}

func (fh *stdHandler) debug(format string, args ...interface{}) {
	fh.logger.Debugf(format, args...)
}

func (fh *stdHandler) error(format string, args ...interface{}) {
	fh.logger.Errorf(format, args...)
}
