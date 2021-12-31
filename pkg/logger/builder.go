package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

// Builder is the builder for logger.
type Builder struct {
	// logger instance
	logger *logrus.Logger
}

// NewBuilder is the startpoint to build a logger.
func NewBuilder() *Builder {
	logger := logrus.New()
	return &Builder{
		logger: logger,
	}
}

// Level sets the level for logger.
func (b *Builder) Level(level logrus.Level) *Builder {
	b.logger.SetLevel(level)
	return b
}

// Rotation accepts a rotate method(lumberjack.Logger, etc) to logger.
func (b *Builder) Rotation(output io.Writer) *Builder {
	b.logger.SetOutput(output)
	return b
}

// Formatter adds a formatter for logger.
func (b *Builder) Formatter(formatter logrus.Formatter) *Builder {
	b.logger.SetFormatter(formatter)
	return b
}

// Complete is the endpoint to finish the builder and return the logger.
func (b *Builder) Complete() *logrus.Logger {
	return b.logger
}

// SetReportCaller whether to log caller info (off by default)
func (b *Builder) SetReportCaller(reportCaller bool) *Builder {
	b.logger.SetReportCaller(reportCaller)
	return b
}
