package glog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	DebugLevel string = "debug"

	InfoLevel string = "info"

	ErrorLevel string = "error"
)

var mapLogs map[string]*Logger = make(map[string]*Logger)
var mu sync.Mutex

type Logger struct {
	atomicLevel zap.AtomicLevel
	log         *zap.Logger
	sugar       *zap.SugaredLogger
}

func getLevelByTag(tag string) zapcore.Level {
	var level zapcore.Level
	switch tag {
	case DebugLevel:
		level = zap.DebugLevel
	case InfoLevel:
		level = zap.InfoLevel
	case ErrorLevel:
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	return level
}
func pathExists(path string) bool {
	path = joinCurrentPath(path)
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
func joinCurrentPath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	if path.IsAbs(p) {
		return p
	}
	return path.Join(getCurrentPath(), p)
}
func getCurrentPath() string {
	dir := filepath.Dir(os.Args[0])
	dir, _ = filepath.Abs(dir)
	return strings.Replace(dir, "\\", "/", -1)
}
func getFilePath(config *LogConfig, args ...string) string {
	parts := make([]string, 0)
	parts = append(parts, joinCurrentPath(config.Path))
	parts = append(parts, "/")
	parts = append(parts, time.Now().Format("20060102"))
	if args != nil && len(args) > 0 {
		parts = append(parts, args...)
	}
	parts = append(parts, ".log")
	return strings.Join(parts, "")
}
func getInstance(key string) *Logger {
	mu.Lock()
	defer mu.Unlock()
	if mapLogs[key] == nil {
		mapLogs[key] = createLogger(key)
	}
	return mapLogs[key]
}
func GetLogger(key string) *Logger {
	return getInstance(key)
}

func createLogger(args ...string) *Logger {
	envConfig := readConfig()
	fileLogger := lumberjack.Logger{
		Filename:   getFilePath(envConfig, args...),
		MaxSize:    10, //MB
		MaxAge:     1,
		MaxBackups: 180,
		LocalTime:  true,
		Compress:   true,
	}
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logLevel := zap.NewAtomicLevel()
	logLevel.SetLevel(getLevelByTag(envConfig.Level))

	encodeTime := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	encoderFileConfig := zap.NewProductionEncoderConfig()
	encoderFileConfig.EncodeTime = encodeTime
	encoderFile := zapcore.NewJSONEncoder(encoderFileConfig)

	encoderConsoleConfig := zap.NewDevelopmentEncoderConfig()
	encoderConsoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConsoleConfig.EncodeTime = encodeTime
	encoderConsole := zapcore.NewConsoleEncoder(encoderConsoleConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(encoderConsole, zapcore.AddSync(os.Stdout), logLevel), //打印到控制台
		zapcore.NewCore(encoderFile, zapcore.AddSync(&fileLogger), logLevel),
	)
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
	l := Logger{
		atomicLevel: logLevel,
		log:         log,
	}
	l.sugar = l.log.Sugar()
	return &l
}
func (l *Logger) Print(v ...interface{}) {
	if v != nil && len(v) > 0 && v[0] == "sql" {
		l.sqlLog(v...)
	} else {
		l.sugar.Debug(v)
	}
}
func (l *Logger) CheckAndPrintError(flag string, err error) {
	if err != nil {
		l.Print(flag, err)
	}
}
func (l *Logger) SetLevel(tag string) {
	l.atomicLevel.SetLevel(getLevelByTag(tag))
}

//nor

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (s *Logger) Debug(msg string, fields ...Field) {
	s.log.Error(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (s *Logger) Info(msg string, fields ...Field) {
	s.log.Error(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (s *Logger) Warn(msg string, fields ...Field) {
	s.log.Error(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (s *Logger) Error(msg interface{}, fields ...Field) error {
	if ev, ok := msg.(error); ok {
		s.log.Error(ev.Error(), fields...)
		return ev
	} else if ev, ok := msg.(string); ok {
		s.log.Error(ev, fields...)
		return fmt.Errorf(ev)
	}
	s.log.Error(fmt.Sprint(msg), fields...)
	return fmt.Errorf(fmt.Sprint(msg))
}
func (s *Logger) Fatal(msg string, fields ...Field) {
	s.log.Error(msg, fields...)
}

//f

// Debugf uses fmt.Sprintf to log a templated message.
func (s *Logger) Debugf(template string, args ...interface{}) {
	s.sugar.Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (s *Logger) Infof(template string, args ...interface{}) {
	s.sugar.Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (s *Logger) Warnf(template string, args ...interface{}) {
	s.sugar.Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (s *Logger) Errorf(template string, args ...interface{}) error {
	s.sugar.Errorf(template, args...)
	return fmt.Errorf(template, args...)
}
func (s *Logger) Fatalf(template string, args ...interface{}) {
	s.sugar.Fatalf(template, args...)
}

// w
// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func (s *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	s.sugar.Debugw(msg, keysAndValues...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (s *Logger) Infow(msg string, keysAndValues ...interface{}) {
	s.sugar.Infow(msg, keysAndValues...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (s *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	s.sugar.Warnw(msg, keysAndValues...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (s *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	s.sugar.Errorw(msg, keysAndValues...)
}
