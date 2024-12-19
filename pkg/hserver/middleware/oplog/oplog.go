package oplog

import (
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
)

var (
	defaultLogger *Logger
	once          sync.Once
)

// Logger 操作日志记录器
type Logger struct {
	writer LogWriter
}

// Init 初始化操作日志记录器
func Init(writer LogWriter) {
	once.Do(func() {
		defaultLogger = &Logger{
			writer: writer,
		}
	})
}

// GetLogger 获取默认的日志记录器
func GetLogger() *Logger {
	return defaultLogger
}

// Record 记录操作日志
func Record(opt LogOption) app.HandlerFunc {
	if defaultLogger == nil {
		panic("oplog logger not initialized")
	}
	return defaultLogger.Record(opt)
}
