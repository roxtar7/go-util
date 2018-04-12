package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Level int

const (
	DEBUG Level = iota + 1
	INFO
	WARN
	ERROR
	FATAL

	Debug = "Debug"
	Info  = "Info"
	Warn  = "Warn"
	Error = "Error"
	Fatal = "Fatal"

	SHORT = iota
	FULL

	FilePathShort = "Short"
	FilePathFull  = "Full"
)

type Logger struct {
	l            *log.Logger
	level        Level
	filePathSize int
	ref          string
	config       *Config
}

func New(config *Config) *Logger {
	l := &Logger{}
	l.config = config
	l.Init(config)
	return l
}

func (l *Logger) Init(config *Config) error {
	l.ref = config.Reference
	l.level = config.Level
	l.filePathSize = config.FilePathSize
	if l.ref == "" {
		refUUID, err := uuid.NewV4()
		if err != nil {
			l.Error("Unable to generate new UUID. Err: ", err)
			return err
		}

		l.ref = refUUID.String()
	}
	l.l = log.New(os.Stdout, fmt.Sprintf("%v [%s] [ %s ] ", time.Now().UTC(), l.config.AppName, l.ref), 0)
	return nil
}

func (l *Logger) GetRef() string {
	return l.ref
}

func (l *Logger) Debug(v ...interface{}) {
	if l.level > DEBUG {
		return
	}

	l.l.Println(l.formatLog("DEBUG", v...)...)
}

func (l *Logger) Info(v ...interface{}) {
	if l.level > INFO {
		return
	}

	l.l.Println(l.formatLog("INFO", v...)...)
}

func (l *Logger) Warn(v ...interface{}) {
	if l.level > WARN {
		return
	}

	l.l.Println(l.formatLog("WARN", v...)...)
}

func (l *Logger) Error(v ...interface{}) {
	if l.level > ERROR {
		return
	}

	l.l.Println(l.formatLog("ERROR", v...)...)
}

func (l *Logger) Fatal(v ...interface{}) {
	if l.level > FATAL {
		return
	}

	l.l.Println(l.formatLog("FATAL", v...)...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level > DEBUG {
		return
	}

	format, v = l.formatLogf("DEBUG", format, v...)
	l.l.Printf(format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level > INFO {
		return
	}

	format, v = l.formatLogf("INFO", format, v...)
	l.l.Printf(format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level > WARN {
		return
	}

	format, v = l.formatLogf("WARN", format, v...)
	l.l.Printf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level > ERROR {
		return
	}

	format, v = l.formatLogf("ERROR", format, v...)
	l.l.Printf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.level > FATAL {
		return
	}

	format, v = l.formatLogf("FATAL", format, v...)
	l.l.Printf(format, v...)
}

// Format the log to contain the log levels
func (l *Logger) formatLog(logType string, v ...interface{}) []interface{} {
	var n []interface{}
	n = append(n, "["+logType+"] ")
	_, file, line, _ := runtime.Caller(2)
	// If you want the short path not the full file path, you can uncomment everything below
	if l.filePathSize == SHORT {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
	}

	n = append(n, file+":"+strconv.Itoa(line)+":")
	n = append(n, v...)
	return n
}

// Format the log to contain the log levels
func (l *Logger) formatLogf(logType string, format string, v ...interface{}) (string, []interface{}) {
	var n []interface{}
	prefix := "[%s] "
	n = append(n, logType)
	_, file, line, _ := runtime.Caller(2)
	// If you want the short path not the full file path, you can uncomment everything below
	if l.filePathSize == SHORT {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
	}

	prefix += "%s:%d: "
	format = prefix + format
	n = append(n, file)
	n = append(n, line)
	n = append(n, v...)
	return format, n
}
