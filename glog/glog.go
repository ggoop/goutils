package glog

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/kataras/golog"
)

type Logger struct {
	*golog.Logger
}

var mapValues map[string]*Logger = make(map[string]*Logger)
var mapKeys map[string]bool = make(map[string]bool)
var mu sync.Mutex
var logDir string = "storage/logs"

func getInstance(key string) *Logger {
	if !mapKeys[key] {
		mu.Lock()
		defer mu.Unlock()
		if !mapKeys[key] {
			mapKeys[key] = true
			mapValues[key] = &Logger{golog.New()}
			mapValues[key].SetLogFile(key)
		}
	}
	return mapValues[key]
}
func createPath(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	return err
}
func (l *Logger) SetLogFile(key string) {
	if err := createPath(logDir); err != nil {
		panic(err)
	}
	filename := path.Join(logDir, key+time.Now().Format("2006-01-02")+".log")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	l.AddOutput(f)
}
func SetPath(path string) {
	logDir = path
}
func GetLogger(key string) *Logger {
	return getInstance(key)
}
func SetLevel(levelName string) {
	getInstance("").SetLevel(levelName)
}
func SetTimeFormat(s string) {
	getInstance("").SetTimeFormat(s)
}
func NewLine(newLine bool) {
	getInstance("").NewLine = newLine
}
func CheckAndPrintError(flag string, err error) {
	if err != nil {
		getInstance("").Println(flag, "\n", err)
	}
}
func Print(v ...interface{}) {
	getInstance("").Print(v...)
}
func Printf(format string, args ...interface{}) {
	getInstance("").Printf(format, args...)
}

// Println prints a log message without levels and colors.
// It adds a new line at the end.
func Println(v ...interface{}) {
	getInstance("").Println(v...)
}

// Fatal `os.Exit(1)` exit no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func Fatal(v ...interface{}) {
	getInstance("").Fatal(v...)
}

// Fatalf will `os.Exit(1)` no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func Fatalf(format string, args ...interface{}) {
	getInstance("").Fatalf(format, args...)
}

// Error will print only when logger's Level is error, warn, info or debug.
func Error(v ...interface{}) error {
	getInstance("").Error(v...)
	return fmt.Errorf(fmt.Sprint(v...))
}

// Errorf will print only when logger's Level is error, warn, info or debug.
func Errorf(format string, args ...interface{}) error {
	getInstance("").Errorf(format, args...)
	return fmt.Errorf(format, args...)
}

// Warn will print when logger's Level is warn, info or debug.
func Warn(v ...interface{}) {
	getInstance("").Warn(v...)
}

// Warnf will print when logger's Level is warn, info or debug.
func Warnf(format string, args ...interface{}) {
	getInstance("").Warnf(format, args...)
}

// Info will print when logger's Level is info or debug.
func Info(v ...interface{}) {
	getInstance("").Info(v...)
}

// Infof will print when logger's Level is info or debug.
func Infof(format string, args ...interface{}) {
	getInstance("").Infof(format, args...)
}

// Debug will print when logger's Level is debug.
func Debug(v ...interface{}) {
	getInstance("").Debug(v...)
}

// Debugf will print when logger's Level is debug.
func Debugf(format string, args ...interface{}) {
	getInstance("").Debugf(format, args...)
}
