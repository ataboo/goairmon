package hardware

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/op/go-logging"
)

func NewLoggingAdaptor(echoLogger echo.Logger) *logging.Logger {
	leveledBackend := logging.SetBackend(&LoggingAdaptor{
		echoLogger,
	})

	logger := logging.MustGetLogger("goairmon")
	logger.SetBackend(leveledBackend)

	return logger
}

type LoggingAdaptor struct {
	EchoLogging echo.Logger
}

func (l *LoggingAdaptor) Log(level logging.Level, callDepth int, record *logging.Record) error {
	switch level {
	case logging.DEBUG:
		l.EchoLogging.Debug(record.Message())
	case logging.INFO:
		l.EchoLogging.Info(record.Message())
	case logging.NOTICE:
		l.EchoLogging.Info(record.Message())
	case logging.WARNING:
		l.EchoLogging.Warn(record.Message())
	case logging.ERROR:
		l.EchoLogging.Error(record.Message())
	case logging.CRITICAL:
		l.EchoLogging.Error(record.Message())
	default:
		return fmt.Errorf("unsupported log level %s", level)
	}

	return nil
}
