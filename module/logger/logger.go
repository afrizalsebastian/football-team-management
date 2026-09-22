package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	ContextKeyMessageId string     = "message_id"
	LogLevelFatal       slog.Level = -1
)

var (
	once     sync.Once
	instance *logger
	logLevel = new(slog.LevelVar)
)

type logger struct {
	*slog.Logger
	attrs []slog.Attr
}

type logBuilder struct {
	logger *logger
	ctx    context.Context
	attrs  []slog.Attr
	msg    string
	level  slog.Level
}

func LoggerNew() *logger {
	once.Do(func() {
		logLevel.Set(slog.LevelDebug)
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.Attr{
						Key:   "timestamp",
						Value: slog.StringValue(time.Now().Format(time.RFC3339Nano)),
					}
				}
				if a.Key == slog.LevelKey {
					lvl := a.Value.Any().(slog.Level)
					label := lvl.String()
					a.Value = slog.StringValue(label)
				}
				return a
			},
		})
		instance = &logger{Logger: slog.New(handler)}
	})
	return instance
}

func (l *logger) createBuilder(msg string, level slog.Level) *logBuilder {
	_, file, line, _ := runtime.Caller(2)
	return &logBuilder{
		logger: l,
		ctx:    context.Background(),
		msg:    msg,
		level:  level,
		attrs: []slog.Attr{
			slog.String("file", filepath.Base(file)),
			slog.Int("line", line),
		},
	}
}

func (l *logger) Debug(msg string) *logBuilder {
	return l.createBuilder(msg, slog.LevelDebug)
}

func (l *logger) Info(msg string) *logBuilder {
	return l.createBuilder(msg, slog.LevelInfo)
}

func (l *logger) Warn(msg string) *logBuilder {
	return l.createBuilder(msg, slog.LevelWarn)
}

func (l *logger) Error(msg string) *logBuilder {
	return l.createBuilder(msg, slog.LevelError)
}

func (l *logger) Fatal(msg string) *logBuilder {
	return l.createBuilder(msg, LogLevelFatal)
}

func (l *logger) Debugf(format string, args ...any) *logBuilder {
	return l.createBuilder(fmt.Sprintf(format, args...), slog.LevelDebug)
}

func (l *logger) Infof(format string, args ...any) *logBuilder {
	return l.createBuilder(fmt.Sprintf(format, args...), slog.LevelInfo)
}

func (l *logger) Warnf(format string, args ...any) *logBuilder {
	return l.createBuilder(fmt.Sprintf(format, args...), slog.LevelWarn)
}

func (l *logger) Errorf(format string, args ...any) *logBuilder {
	return l.createBuilder(fmt.Sprintf(format, args...), slog.LevelError)
}

func (l *logger) Fatalf(format string, args ...any) *logBuilder {
	return l.createBuilder(fmt.Sprintf(format, args...), LogLevelFatal)
}

func (l *logger) WithContext(ctx context.Context) *logBuilder {
	return &logBuilder{
		logger: l,
		ctx:    ctx,
	}
}

func (lb *logBuilder) getDetailAttributes() []slog.Attr {
	programCounter, file, line, _ := runtime.Caller(2)
	lastPath := strings.Split(filepath.Base(runtime.FuncForPC(programCounter).Name()), ".")

	attrs := []slog.Attr{
		slog.String("file", filepath.Base(file)),
		slog.Int("line", line),
		slog.String("package", lastPath[0]),
	}

	if len(lastPath) > 2 {
		attrs = append(attrs, slog.String("struct", lastPath[1][1:len(lastPath[1])-1]))
	}

	attrs = append(attrs, slog.String("function", lastPath[len(lastPath)-1]))
	return attrs
}

func (lb *logBuilder) Debug(msg string) *logBuilder {
	lb.attrs = append(lb.attrs, lb.getDetailAttributes()...)
	lb.level = slog.LevelDebug
	lb.msg = msg
	return lb
}

func (lb *logBuilder) Info(msg string) *logBuilder {
	lb.attrs = append(lb.attrs, lb.getDetailAttributes()...)
	lb.level = slog.LevelInfo
	lb.msg = msg
	return lb
}

func (lb *logBuilder) Warn(msg string) *logBuilder {
	lb.attrs = append(lb.attrs, lb.getDetailAttributes()...)
	lb.level = slog.LevelWarn
	lb.msg = msg
	return lb
}

func (lb *logBuilder) Error(msg string) *logBuilder {
	lb.attrs = append(lb.attrs, lb.getDetailAttributes()...)
	lb.level = slog.LevelError
	lb.msg = msg
	return lb
}

func (lb *logBuilder) Fatal(msg string) *logBuilder {
	lb.attrs = append(lb.attrs, lb.getDetailAttributes()...)
	lb.level = LogLevelFatal
	lb.msg = msg
	return lb
}

func (lb *logBuilder) Debugf(format string, args ...any) *logBuilder {
	lb.attrs = append(lb.attrs, lb.getDetailAttributes()...)
	lb.level = slog.LevelDebug
	lb.msg = fmt.Sprintf(format, args...)
	return lb
}

func (lb *logBuilder) Infof(format string, args ...any) *logBuilder {
	lb.attrs = append(lb.attrs, lb.getDetailAttributes()...)
	lb.level = slog.LevelInfo
	lb.msg = fmt.Sprintf(format, args...)
	return lb
}

func (lb *logBuilder) Warnf(format string, args ...any) *logBuilder {
	lb.attrs = append(lb.attrs, lb.getDetailAttributes()...)
	lb.level = slog.LevelWarn
	lb.msg = fmt.Sprintf(format, args...)
	return lb
}

func (lb *logBuilder) Errorf(format string, args ...any) *logBuilder {
	lb.attrs = append(lb.attrs, lb.getDetailAttributes()...)
	lb.level = slog.LevelError
	lb.msg = fmt.Sprintf(format, args...)
	return lb
}

func (lb *logBuilder) Fatalf(format string, args ...any) *logBuilder {
	lb.attrs = append(lb.attrs, lb.getDetailAttributes()...)
	lb.level = LogLevelFatal
	lb.msg = fmt.Sprintf(format, args...)
	return lb
}

func (lb *logBuilder) Attr(key string, value any) *logBuilder {
	lb.attrs = append(lb.attrs, slog.Any(key, value))
	return lb
}

func (lb *logBuilder) Attrs(args ...string) *logBuilder {
	for i := 0; i < len(args)-1; i += 2 {
		lb.attrs = append(lb.attrs, slog.Any(args[i], args[i+1]))
	}
	return lb
}

// for showing the msg in console
func (lb *logBuilder) Msg(msg ...any) {
	allAttrs := append(lb.logger.attrs, lb.getIdentificationAttributes()...)
	allAttrs = append(allAttrs, lb.attrs...)

	logArgs := convertAttrs(allAttrs)
	logArgs = append(logArgs, msg...)

	switch lb.level {
	case slog.LevelDebug:
		lb.logger.Logger.With(logArgs...).Debug(lb.msg)
	case slog.LevelInfo:
		lb.logger.Logger.With(logArgs...).Info(lb.msg)
	case slog.LevelWarn:
		lb.logger.Logger.With(logArgs...).Warn(lb.msg)
	case slog.LevelError:
		lb.logger.Logger.With(logArgs...).Error(lb.msg)
	case LogLevelFatal:
		lb.logger.Logger.With(logArgs...).Error(lb.msg)
		os.Exit(1)
	default:
		lb.logger.Logger.With(logArgs...).Info(lb.msg)
	}

	lb.attrs = []slog.Attr{}
}

func (lb *logBuilder) getIdentificationAttributes() []slog.Attr {
	var attrs []slog.Attr

	if val, ok := lb.ctx.Value(ContextKeyMessageId).(string); ok {
		attrs = append(attrs, slog.String(ContextKeyMessageId, val))
	}

	// for other identification attributes
	return attrs
}

func convertAttrs(attrs []slog.Attr) []any {
	args := make([]any, 0, len(attrs)*2)
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value)
	}
	return args
}
