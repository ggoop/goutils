package glog

// Default is the package-level ready-to-use logger,
// level had set to "info", is changeable.
var Default = createLogger()

// Reset re-sets the default logger to an empty one.
func Reset() {
	Default = createLogger()
}

// Print prints a log message without levels and colors.
func Print(v ...interface{}) {
	Default.Print(v...)
}

// Println prints a log message without levels and colors.
// It adds a new line at the end.
func Println(v ...interface{}) {
	Default.sugar.Info(v...)
}
func CheckAndPrintError(flag string, err error) {
	Default.CheckAndPrintError(flag, err)
}

// Fatal `os.Exit(1)` exit no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func Fatal(msg string, fields ...Field) {
	Default.Fatal(msg, fields...)
}

// Fatalf will `os.Exit(1)` no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func Fatalf(format string, args ...interface{}) {
	Default.Fatalf(format, args...)
}

// Error will print only when logger's Level is error, warn, info or debug.
func Error(msg interface{}, fields ...Field) error {
	return Default.Error(msg, fields...)

}

// Errorf will print only when logger's Level is error, warn, info or debug.
func Errorf(format string, args ...interface{}) error {
	return Default.Errorf(format, args...)
}

// Warn will print when logger's Level is warn, info or debug.
func Warn(msg string, fields ...Field) {
	Default.Warn(msg, fields...)
}

// Warnf will print when logger's Level is warn, info or debug.
func Warnf(format string, args ...interface{}) {
	Default.Warnf(format, args...)
}

// Info will print when logger's Level is info or debug.
func Info(msg string, fields ...Field) {
	Default.Info(msg, fields...)
}

// Infof will print when logger's Level is info or debug.
func Infof(format string, args ...interface{}) {
	Default.Infof(format, args...)
}

// Debug will print when logger's Level is debug.
func Debug(msg string, fields ...Field) {
	Default.Debug(msg, fields...)
}

// Debugf will print when logger's Level is debug.
func Debugf(format string, args ...interface{}) {
	Default.Debugf(format, args...)
}
