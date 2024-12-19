package oplog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// LogWriter 日志写入器接口
type LogWriter interface {
	Write(ctx context.Context, log *OperationLog) error
	Close() error
}

// FileWriter 文件写入器
type FileWriter struct {
	logDir string
	file   *os.File
	date   string
	mutex  sync.Mutex
}

func NewFileWriter(logDir string) *FileWriter {
	return &FileWriter{
		logDir: logDir,
	}
}

func (w *FileWriter) Write(ctx context.Context, log *OperationLog) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	date := log.CreatedAt.Format("2006-01-02")
	if w.file == nil || w.date != date {
		if w.file != nil {
			w.file.Close()
		}
		if err := w.rotateFile(date); err != nil {
			return err
		}
	}

	// 格式化日志内容
	logContent := fmt.Sprintf("[%s] UserID:%s Username:%s TenantID:%s Module:%s Action:%s Method:%s Path:%s Query:%s IP:%s UA:%s Status:%d Duration:%dms Body:%s Error:%s\n",
		log.CreatedAt.Format("2006-01-02 15:04:05"),
		log.UserID,
		log.Username,
		log.TenantID,
		log.Module,
		log.Action,
		log.Method,
		log.Path,
		log.Query,
		log.IP,
		log.UserAgent,
		log.Status,
		log.Duration,
		log.Body,
		log.Error,
	)

	_, err := w.file.WriteString(logContent)
	return err
}

func (w *FileWriter) rotateFile(date string) error {
	logFile := filepath.Join(w.logDir, fmt.Sprintf("oplog_%s.log", date))
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	w.file = f
	w.date = date
	return nil
}

func (w *FileWriter) Close() error {
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

// DBWriter 数据库写入器
type DBWriter struct {
	db interface{} // 使用您的数据库接口
}

func NewDBWriter(db interface{}) *DBWriter {
	return &DBWriter{db: db}
}

func (w *DBWriter) Write(ctx context.Context, log *OperationLog) error {
	// 根据月份确定表名
	tableName := fmt.Sprintf("sys_operation_log_%s", log.CreatedAt.Format("200601"))

	// 确保表存在
	if err := w.ensureTable(ctx, tableName); err != nil {
		return err
	}

	// 写入日志
	// 实现数据库写入逻辑
	return nil
}

func (w *DBWriter) ensureTable(ctx context.Context, tableName string) error {
	// 实现建表逻辑
	return nil
}

func (w *DBWriter) Close() error {
	return nil
}

// MultiWriter 多重写入器
type MultiWriter struct {
	writers []LogWriter
}

func NewMultiWriter(writers ...LogWriter) *MultiWriter {
	return &MultiWriter{writers: writers}
}

func (w *MultiWriter) Write(ctx context.Context, log *OperationLog) error {
	for _, writer := range w.writers {
		if err := writer.Write(ctx, log); err != nil {
			hlog.CtxErrorf(ctx, "write operation log error: %v", err)
		}
	}
	return nil
}

func (w *MultiWriter) Close() error {
	for _, writer := range w.writers {
		if err := writer.Close(); err != nil {
			hlog.Errorf("close log writer error: %v", err)
		}
	}
	return nil
}
