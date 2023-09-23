package logger

import (
	"fmt"
	"os"
	"time"

	"log/slog"

	_ "net/http/pprof"
)

type Logger struct {
	title     string
	log       *slog.Logger
	startTime time.Time
}

func InitLogger(title string) *Logger {
	l := &Logger{
		title:     title,
		startTime: time.Now(),
	}

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}).WithAttrs([]slog.Attr{slog.String("service", l.title)})
	l.log = slog.New(textHandler)
	return l
}

func (l Logger) GetSysLogger() *slog.Logger {
	return l.log
}

func (l Logger) Info(msg string) {
	l.log.Info(msg,
		slog.Duration("Uptime", time.Since(l.startTime)))
}

func (l Logger) Infow(msg string, args ...any) {
	l.log.Info(msg, args,
		slog.Duration("Uptime", time.Since(l.startTime)))
}

func (l Logger) Infof(msg string, args ...any) {
	l.log.Info(fmt.Sprintf(msg, args...),
		slog.Duration("Uptime", time.Since(l.startTime)))
}

func (l Logger) Warning(msg string) {
	l.log.Warn(msg,
		slog.Duration("Uptime", time.Since(l.startTime)))
}

func (l Logger) Warningw(msg string, args ...any) {
	l.log.Warn(msg, args,
		slog.Duration("Uptime", time.Since(l.startTime)))
}

func (l Logger) Warningf(msg string, args ...any) {
	l.log.Warn(fmt.Sprintf(msg, args...),
		slog.Duration("Uptime", time.Since(l.startTime)))
}

func (l Logger) Error(msg string) {
	l.log.Error(msg,
		slog.Duration("Uptime", time.Since(l.startTime)))
}

func (l Logger) Errorw(msg string, args ...any) {
	l.log.Error(msg, args,
		slog.Duration("Uptime", time.Since(l.startTime)))
}

func (l Logger) Errorf(msg string, args ...any) {
	l.log.Error(fmt.Sprintf(msg, args...),
		slog.Duration("Uptime", time.Since(l.startTime)))
}

func (l Logger) Debug(msg string) {
	l.log.Debug(msg,
		slog.Duration("Uptime", time.Since(l.startTime)))
}

func (l Logger) Debugw(msg string, args ...any) {
	l.log.Debug(msg, args,
		slog.Duration("Uptime", time.Since(l.startTime)))
}

func (l Logger) Debugf(msg string, args ...any) {
	l.log.Debug(fmt.Sprintf(msg, args...),
		slog.Duration("Uptime", time.Since(l.startTime)))
}
