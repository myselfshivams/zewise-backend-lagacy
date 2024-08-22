/*
Package loggers - NekoBlog backend server loggers.
This file is for logger factory and logger formatter.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package loggers

import (
	"bytes"
	"strings"

	"github.com/sirupsen/logrus"
)

// LoggerFormatter 日志格式化器1
type LoggerFormatter struct {
}

// Format 格式化日志
func (formatter *LoggerFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 时间、日志级别、日志内容
	var sb strings.Builder
	sb.WriteString("[")
	sb.WriteString(entry.Time.Format("15:04:05"))
	sb.WriteString("][")

	// 日志级别颜色
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

	// 日志级别
	sb.WriteString(strings.ToUpper(entry.Level.String()))

	// 还原颜色
	sb.WriteString("\x1b[0m]")
	prefix := sb.String()

	// 写入日志内容
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

	// 返回日志内容
	return buffer.Bytes(), nil
}

// 创建 logger 工厂函数
func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)            // 设置报告调用者
	logger.SetFormatter(&LoggerFormatter{}) // 设置日志格式化器

	return logger
}
