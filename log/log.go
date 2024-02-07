// Package log provides and implements a simple leveled logging interface.
// It defines a type, Logger, with methods for formatting output.
// Helper functions provided are easier to use than creating a Logger manually.
package log

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

var defaultLogger *logrus.Entry = logrus.StandardLogger().WithField("go.version", runtime.Version())

// Logger provides a leveled-logging interface.
type Logger interface {
	// Print calls Print on logger implementation to print to the logger.
	// Arguments are handled in the manner of fmt.Print.
	Print(args ...interface{})

	// Printf calls Printf on logger implementation to print to the logger.
	// Arguments are handled in the manner of fmt.Printf.
	Printf(format string, args ...interface{})

	// Println calls Prinln on logger implementation to print to the logger.
	// Arguments are handled in the manner of fmt.Println.
	Println(args ...interface{})

	// Fatal is equivalent calling Print() on logger implementation followed by a call to os.Exit(1).
	Fatal(args ...interface{})

	// Fatalf is equivalent calling Printf() on logger implementation followed by a call to os.Exit(1).
	Fatalf(format string, args ...interface{})

	// Fatalln is equivalent calling Println() on logger implementation followed by a call to os.Exit(1).
	Fatalln(args ...interface{})

	// Panic is equivalent calling Print() on logger implementation followed by a call to panic().
	Panic(args ...interface{})

	// Panicf is equivalent calling Printf() on logger implementation followed by a call to panic().
	Panicf(format string, args ...interface{})

	// Panicln is equivalent Println() on logger implementation followed by a call to panic().
	Panicln(args ...interface{})

	// Debug is equivalent calling Print() on logger implementation followed by a adding Debug level.
	Debug(args ...interface{})

	// Debugf is equivalent calling Printf() on logger implementation followed by a adding Debug level.
	Debugf(format string, args ...interface{})

	// Debugln is equivalent calling Println() on logger implementation followed by a adding Debug level.
	Debugln(args ...interface{})

	// Error is equivalent calling Print() on logger implementation followed by a adding Error level.
	Error(args ...interface{})

	// Errorf is equivalent calling Printf() on logger implementation followed by a adding Error level.
	Errorf(format string, args ...interface{})

	// Errorln is equivalent calling Println() on logger implementation followed by a adding Error level.
	Errorln(args ...interface{})

	// Info is equivalent calling Print() on logger implementation followed by a adding Info level.
	Info(args ...interface{})

	// Infof is equivalent calling Printf() on logger implementation followed by a adding Info level.
	Infof(format string, args ...interface{})

	// Infoln is equivalent calling Println() on logger implementation followed by a adding Info level.
	Infoln(args ...interface{})

	// Warn is equivalent calling Print() on logger implementation followed by a adding Warn level.
	Warn(args ...interface{})

	// Warnf is equivalent calling Printf() on logger implementation followed by a adding Warn level.
	Warnf(format string, args ...interface{})

	// Warnln is equivalent calling Println() on logger implementation followed by a adding Warn level.
	Warnln(args ...interface{})

	// WithError wraps an error.
	WithError(err error) *logrus.Entry
}

// Fields is used as argument in WithFields method/func
type Fields = logrus.Fields

// Default is a default logger instance
func Default() Logger {
	return defaultLogger
}

// Print calls defaultLogger.Print to print to the logger.
func Print(args ...interface{}) { defaultLogger.Print(args...) }

// Printf calls defaultLogger.Printf to print to the logger.
func Printf(format string, args ...interface{}) { defaultLogger.Printf(format, args...) }

// Println calls defaultLogger.Println to print to the logger.
func Println(args ...interface{}) { defaultLogger.Println(args...) }

// Fatal calls defaultLogger.Fatal to print to the logger.
func Fatal(args ...interface{}) { defaultLogger.Fatal(args...) }

// Fatalf calls defaultLogger.Fatalf to print to the logger.
func Fatalf(format string, args ...interface{}) { defaultLogger.Fatalf(format, args...) }

// Fatalln calls defaultLogger.Fatalln to print to the logger.
func Fatalln(args ...interface{}) { defaultLogger.Fatalln(args...) }

// Panic calls defaultLogger.Panic to print to the logger.
func Panic(args ...interface{}) { defaultLogger.Panic(args...) }

// Panicf calls defaultLogger.Panicf to print to the logger.
func Panicf(format string, args ...interface{}) { defaultLogger.Panicf(format, args...) }

// Panicln calls defaultLogger.Panicln to print to the logger.
func Panicln(args ...interface{}) { defaultLogger.Panicln(args...) }

// Debug calls defaultLogger.Debug to print to the logger.
func Debug(args ...interface{}) { defaultLogger.Debug(args...) }

// Debugf calls defaultLogger.Debugf to print to the logger.
func Debugf(format string, args ...interface{}) { defaultLogger.Debugf(format, args...) }

// Debugln calls defaultLogger.Debugln to print to the logger.
func Debugln(args ...interface{}) { defaultLogger.Debugln(args...) }

// Error calls defaultLogger.Error to print to the logger.
func Error(args ...interface{}) { defaultLogger.Error(args...) }

// Errorf calls defaultLogger.Errorf to print to the logger.
func Errorf(format string, args ...interface{}) { defaultLogger.Errorf(format, args...) }

// Errorln calls defaultLogger.Errorln to print to the logger.
func Errorln(args ...interface{}) { defaultLogger.Errorln(args...) }

// Info calls defaultLogger.Info to print to the logger.
func Info(args ...interface{}) { defaultLogger.Info(args...) }

// Infof calls defaultLogger.Infof to print to the logger.
func Infof(format string, args ...interface{}) { defaultLogger.Infof(format, args...) }

// Infoln calls defaultLogger.Infoln to print to the logger.
func Infoln(args ...interface{}) { defaultLogger.Infoln(args...) }

// Warn calls defaultLogger.Warn to print to the logger.
func Warn(args ...interface{}) { defaultLogger.Warn(args...) }

// Warnf calls defaultLogger.Warnf to print to the logger.
func Warnf(format string, args ...interface{}) { defaultLogger.Warnf(format, args...) }

// Warnln calls defaultLogger.Warnln to print to the logger.
func Warnln(args ...interface{}) { defaultLogger.Warnln(args...) }

// An entry is the final or intermediate logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Trace, Debug,
// Info, Warn, Error, Fatal or Panic is called on it. These objects can be
// reused and passed around as much as you wish to avoid field duplication.
type Entry struct {
	*logrus.Entry
}

// Add an error as single field (using the key defined in ErrorKey) to the Entry.
func WithError(err error) *Entry {
	var entry Entry

	entry.Entry = defaultLogger.WithError(err)
	return &entry
}

// Add a map of fields to the Entry.
func WithFields(fields Fields) *Entry {
	var entry Entry

	entry.Entry = defaultLogger.WithFields(logrus.Fields(fields))
	return &entry
}

// Add a map of fields to the Entry.
func (e *Entry) WithFields(fields Fields) *Entry {
	e.Entry = e.Entry.WithFields(logrus.Fields(fields))
	return e
}
