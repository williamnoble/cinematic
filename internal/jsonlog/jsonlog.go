package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Level to represent the severity of a log entry
type Level int8

// Constant which represent the severity level.
const (
	LevelInfo  Level = iota //0
	LevelError              //1
	LevelFatal              //2
	LevelOff                //3
)

// String prints a human readable string.
func (l Level) String() string {
	switch l {
	case LevelOff:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// Logger defines a custom logger. It holds an output destination, the minimum
// severity level and a mutex to avoid race conditions when writing data.
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

// New creates a new Logger, it writers to specified writer e.g. Stdout and takes
// in a minLevel parameter which describes the minimum severity of logging.
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

// PrintInfo prints out general information message with unlimited number of props
// example: logger.PrintInfo("database connection pool established", nil)
func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

// PrintError writes a standard error to given writer
func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

// PrintFatal will print Fatal errors and then suspend exucution via calling os.Exit(1)
func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

// print is a method function of Logger which returns a json object containing
//the level of severity, a message, a map of properties.
func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string
		Time       string
		Message    string
		Properties map[string]string
		Trace      string // Not for LevelInfo
	}{
		Level:      level.String(),
		Time:       time.Now().Format(time.RFC3339), // return time at UTC
		Message:    message,
		Properties: properties,
	}

	// Ignoring Information errors, we include a stack trace for more severe errors.
	if level >= LevelError { // level > 1
		aux.Trace = string(debug.Stack())
	}

	var line []byte // error line

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message" + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}

func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
