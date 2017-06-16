// Package logger aims at providing an advanced implementation of the logging.
package logger

import (
	"log"
	"os"
	"strconv"
)

// Logger is a structure in charge of logging messages.
type Logger struct {
	DebugEnabled bool
	stdLogger    *log.Logger
	debugLogger  *log.Logger
}

var (
	logger *Logger
)

func init() {
	logger = new(Logger)
	logger.DebugEnabled, _ = strconv.ParseBool(os.Getenv("LOG_DEBUG"))
	logger.stdLogger = log.New(os.Stdout, "", log.LstdFlags)
	logger.debugLogger = log.New(os.Stdout, "DEBUG: ", log.LstdFlags)
}

// GetLogger returns a pointer to the existing logger.
func GetLogger() *Logger {
	return logger
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (logger *Logger) Printf(format string, values ...interface{}) {
	logger.stdLogger.Printf(format, values...)
}

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (logger *Logger) Println(values ...interface{}) {
	logger.stdLogger.Println(values...)
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (logger *Logger) Print(values ...interface{}) {
	logger.stdLogger.Print(values...)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (logger *Logger) Fatalf(format string, values ...interface{}) {
	logger.stdLogger.Fatalf(format, values...)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (logger *Logger) Fatalln(values ...interface{}) {
	logger.stdLogger.Fatalln(values...)
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (logger *Logger) Fatal(values ...interface{}) {
	logger.stdLogger.Fatal(values...)
}

// Debugf calls l.Output to print to the logger if the debug mode is enabled.
// Arguments are handled in the manner of fmt.Printf.
func (logger *Logger) Debugf(format string, values ...interface{}) {
	if logger.DebugEnabled {
		logger.debugLogger.Printf(format, values...)
	}
}

// Debugln calls l.Output to print to the logger if the debug mode is enabled.
// Arguments are handled in the manner of fmt.Println.
func (logger *Logger) Debugln(values ...interface{}) {
	if logger.DebugEnabled {
		logger.debugLogger.Println(values...)
	}
}

// Debug calls l.Output to print to the logger if the debug mode is enabled.
// Arguments are handled in the manner of fmt.Print.
func (logger *Logger) Debug(values ...interface{}) {
	if logger.DebugEnabled {
		logger.debugLogger.Print(values...)
	}
}
