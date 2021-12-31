package logger

import (
	"context"
	"fmt"
	"io"
	"runtime/debug"
	"time"

	cfg "github.com/dollarkillerx/common/pkg/config"
	"github.com/sirupsen/logrus"
)

// RimeLogger ...
type RimeLogger struct {
	*logrus.Logger
}

// NewLogger is the startpoint to build a Logger.
func NewLogger() *RimeLogger {
	logger := logrus.New()
	return &RimeLogger{
		logger,
	}
}

func NewRimeLogger(loggerConfig cfg.LoggerConfig) *RimeLogger {
	return NewLogger().
		Level(loggerConfig.Level.Level()).
		Formatter(DefaultFormatter()).
		Rotation(DefaultRotation(&loggerConfig)).
		SetLogReportCaller(true).
		Complete()
}

// Level sets the level for logger.
func (l *RimeLogger) Level(level logrus.Level) *RimeLogger {
	l.SetLevel(level)
	return l
}

// Rotation accepts a rotate method(lumberjack.Logger, etc) to logger.
func (l *RimeLogger) Rotation(output io.Writer) *RimeLogger {
	l.SetOutput(output)
	return l
}

// Formatter adds a formatter for logger.
func (l *RimeLogger) Formatter(formatter logrus.Formatter) *RimeLogger {
	l.SetFormatter(formatter)
	return l
}

// Complete is return the Logger.
func (l *RimeLogger) Complete() *RimeLogger {
	return l
}

// SetLogReportCaller whether to log caller info (off by default)
func (l *RimeLogger) SetLogReportCaller(reportCaller bool) *RimeLogger {
	l.SetReportCaller(reportCaller)
	return l
}

// Recover ...
// 记录 panic 的堆栈信息, 存储 req 信息
func (l *RimeLogger) Recover(ctx context.Context, req interface{}) {
	if r := recover(); r != nil {
		l.Errorf("[ReqID:%s] [%s] [Panic]  [%+v]", "ReqID", "ServerName", string(debug.Stack()))
	}

	dataMap, ok := ctx.Value(dataMapCtxKey).(map[string]interface{})
	if ok {
		dataMap[httpRequestKey] = req
	}
	return
}

// Print ...
// 记录慢 SQL，gorm 自动调用，超过 100ms 认定为慢 SQL
func (l *RimeLogger) Print(values ...interface{}) {
	if len(values) >= 3 && values[0] == "sql" {
		duration := values[2].(time.Duration)
		if duration > time.Millisecond*100 {
			l.Warningf("[Slow Query] [%s] [%s]", values[1], values[2])
		}
	}
}

// ErrorHandler ... 存储 error 的 stack
// Deprecated 新的 errs pkg 不需要在最外层使用这个函数
func (l *RimeLogger) ErrorHandler(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	var logErrorMessage string // 记录到文件的日志信息
	var displayError error     // 返回到前端的错误
	var errID string

	logErrorMessage = fmt.Sprintf("%+v", err) // 记录到文件的日志信息
	displayError = err                        // 返回到前端的错误

	dataMap, ok := ctx.Value(dataMapCtxKey).(map[string]interface{})
	if ok {
		dataMap[errorWithStackKey] = logErrorMessage
		dataMap[errorID] = errID
	}

	return displayError
}
