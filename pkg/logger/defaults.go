package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	cfg "github.com/dollarkillerx/common/pkg/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// DefaultFormatter returns a default json formatter
func DefaultFormatter() logrus.Formatter {
	return &logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Dir(f.File) + "/" + path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}
}

// DefaultRotation returns a default Rotation
func DefaultRotation(loggerConfig *cfg.LoggerConfig) io.Writer {
	if loggerConfig.Filename == "" {
		return os.Stdout
	}

	result := &lumberjack.Logger{
		Filename: loggerConfig.Filename,
		Compress: loggerConfig.Compress,
	}
	if loggerConfig.MaxSize != 0 {
		result.MaxSize = loggerConfig.MaxSize
	}
	if loggerConfig.MaxAge != 0 {
		result.MaxAge = loggerConfig.MaxAge
	}
	if loggerConfig.MaxBackups != 0 {
		result.MaxBackups = loggerConfig.MaxBackups
	}

	return result
}

type CustomizeLogWrite struct {
	wFun  func(p []byte) (n int, err error)
	log   func(p []byte)
	write io.Writer
}

func (c *CustomizeLogWrite) Write(p []byte) (n int, err error) {
	return c.wFun(p)
}

func DefaultCustomizeLog(loggerConfig *cfg.LoggerConfig, log func(p []byte)) *CustomizeLogWrite {
	var customizeLogrWrite CustomizeLogWrite

	if loggerConfig.Filename == "" {
		customizeLogrWrite.write = os.Stdout
	} else {
		result := &lumberjack.Logger{
			Filename: loggerConfig.Filename,
			Compress: loggerConfig.Compress,
		}
		if loggerConfig.MaxSize != 0 {
			result.MaxSize = loggerConfig.MaxSize
		}
		if loggerConfig.MaxAge != 0 {
			result.MaxAge = loggerConfig.MaxAge
		}
		if loggerConfig.MaxBackups != 0 {
			result.MaxBackups = loggerConfig.MaxBackups
		}
		customizeLogrWrite.write = result
	}

	if log != nil {
		customizeLogrWrite.log = log
	}

	customizeLogrWrite.wFun = func(p []byte) (n int, err error) {
		if customizeLogrWrite.log != nil {
			customizeLogrWrite.log(p)
		}

		return customizeLogrWrite.write.Write(p)
	}

	return &customizeLogrWrite
}
