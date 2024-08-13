package loggers

import (
	"bytes"
	"strings"

	"github.com/sirupsen/logrus"
)

type LoggerFormatter struct {
}

func (formatter *LoggerFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var sb strings.Builder
	sb.WriteString("[")
	sb.WriteString(entry.Time.Format("15:04:05"))
	sb.WriteString("][")

	var levelColor string
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = "\x1b[32m"
	case logrus.WarnLevel:
		levelColor = "\x1b[33m"
	case logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel:
		levelColor = "\x1b[31m"
	case logrus.DebugLevel:
		levelColor = "\x1b[90m"
	default:
		levelColor = "\x1b[0m"
	}
	sb.WriteString(levelColor)

	sb.WriteString(strings.ToUpper(entry.Level.String()))

	sb.WriteString("\x1b[0m]")
	prefix := sb.String()

	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = new(bytes.Buffer)
	}

	buffer.WriteString(prefix)
	buffer.WriteRune(' ')
	buffer.WriteString(entry.Message)
	buffer.WriteRune('\n')

	return buffer.Bytes(), nil
}

// 创建 logger 工厂函数
func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(&LoggerFormatter{})

	return logger
}
