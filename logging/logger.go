package logging

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LogLevel int

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
	PANIC
	OFF
)

var root *logger

type Logger interface {
	SetName(name string)
	Debug(msg string)
	Debugf(format string, a ...any)
	Info(msg string)
	Infof(format string, a ...any)
	Warn(msg string)
	Warnf(format string, a ...any)
	Error(msg string)
	Errorf(format string, a ...any)
	Fatal(msg string)
	Fatalf(format string, a ...any)
	Panic(msg string)
	Panicf(format string, a ...any)
}

type Config struct {
	DisableFileLog bool
	DisableStdLog  bool
	LogQueue       int
	LogFile        string
	LogMaxDays     int
	LogMaxBackups  int
	LogMaxSize     int
}

type logMessage struct {
	msg   string
	fatal bool
	panic bool
}

type logger struct {
	LogLevel LogLevel
	lumber   lumberjack.Logger
	format   string
	name     string
	queue    chan logMessage
	quit     chan bool
	template *template.Template
	m        map[string]string
}

func Initialize(config *Config) {
	fmt.Printf("%#v\n", config)
	root = &logger{
		name:   "root",
		format: "{{.time}} :: {{.level}} : {{.logger}} : {{.msg}}",
		queue:  make(chan logMessage, config.LogQueue),
		quit:   make(chan bool),
	}
	if !config.DisableFileLog {
		root.lumber = lumberjack.Logger{
			Filename:   config.LogFile,
			MaxSize:    config.LogMaxSize,
			MaxBackups: config.LogMaxBackups,
			MaxAge:     config.LogMaxDays,
			Compress:   true,
		}
	}

	if !strings.HasSuffix(root.format, "\n") {
		root.format += "\n"
	}
	root.template = template.Must(template.New("").Parse(root.format))

	filter := regexp.MustCompile("\u001B\\[\\d+m")

	go func() {
		for {
			select {
			case <-root.quit:
				err := root.lumber.Close()
				if err != nil {
					fmt.Printf("Error closing log file")
				}
				return
			case msg := <-root.queue:
				if !config.DisableStdLog {
					_, err := os.Stdout.WriteString(msg.msg)
					if err != nil {
						fmt.Printf("Error writing to stdout: %s\n", err.Error())
					}
				}
				if !config.DisableFileLog {
					msg.msg = filter.ReplaceAllString(msg.msg, "")
					_, err := root.lumber.Write([]byte(msg.msg))
					if err != nil {
						fmt.Printf("Error writing logs: %s\n", err.Error())
					}
				}
				if msg.fatal {
					os.Exit(1)
				} else if msg.panic {
					panic(msg.msg)
				}
			}
		}
	}()
}

func GetLogger(name string) Logger {
	l := &logger{
		format:   root.format,
		template: root.template,
		LogLevel: root.LogLevel,
		m:        make(map[string]string),
	}
	l.SetName(name)
	return l
}

func (l *logger) Debug(msg string) {
	if l.LogLevel > DEBUG {
		return
	}
	buffer := bytes.Buffer{}
	l.m["time"] = time.Now().Format(time.RFC3339Nano)
	l.m["msg"] = msg
	l.m["level"] = "DEBUG"
	l.template.Execute(&buffer, l.m)
	root.queue <- logMessage{msg: buffer.String()}
}

func (l *logger) Debugf(format string, a ...any) {
	if l.LogLevel > DEBUG {
		return
	}
	buffer := bytes.Buffer{}
	l.m["time"] = time.Now().Format(time.RFC3339Nano)
	l.m["msg"] = fmt.Sprintf(format, a...)
	l.m["level"] = "DEBUG"
	l.template.Execute(&buffer, l.m)
	root.queue <- logMessage{msg: buffer.String()}
}

func (l *logger) Info(msg string) {
	if l.LogLevel > INFO {
		return
	}
	buffer := bytes.Buffer{}
	l.m["time"] = time.Now().Format(time.RFC3339Nano)
	l.m["msg"] = msg
	l.m["level"] = "INFO "
	l.template.Execute(&buffer, l.m)
	root.queue <- logMessage{msg: buffer.String()}

}

func (l *logger) Infof(format string, a ...any) {
	if l.LogLevel > INFO {
		return
	}
	buffer := bytes.Buffer{}
	l.m["time"] = time.Now().Format(time.RFC3339Nano)
	l.m["msg"] = fmt.Sprintf(format, a...)
	l.m["level"] = "INFO "
	l.template.Execute(&buffer, l.m)
	root.queue <- logMessage{msg: buffer.String()}
}

func (l *logger) Warn(msg string) {
	if l.LogLevel > WARN {
		return
	}
	buffer := bytes.Buffer{}
	l.m["time"] = time.Now().Format(time.RFC3339Nano)
	l.m["msg"] = msg
	l.m["level"] = "WARN "
	l.template.Execute(&buffer, l.m)
	root.queue <- logMessage{msg: buffer.String()}
}

func (l *logger) Warnf(format string, a ...any) {
	if l.LogLevel > WARN {
		return
	}
	buffer := bytes.Buffer{}
	l.m["time"] = time.Now().Format(time.RFC3339Nano)
	l.m["msg"] = fmt.Sprintf(format, a...)
	l.m["level"] = "WARN "
	l.template.Execute(&buffer, l.m)
	root.queue <- logMessage{msg: buffer.String()}
}

func (l *logger) Error(msg string) {
	if l.LogLevel > ERROR {
		return
	}
	buffer := bytes.Buffer{}
	l.m["time"] = time.Now().Format(time.RFC3339Nano)
	l.m["msg"] = msg
	l.m["level"] = "ERROR"
	l.template.Execute(&buffer, l.m)
	root.queue <- logMessage{msg: buffer.String()}
}

func (l *logger) Errorf(format string, a ...any) {
	if l.LogLevel > ERROR {
		return
	}
	buffer := bytes.Buffer{}
	l.m["time"] = time.Now().Format(time.RFC3339Nano)
	l.m["msg"] = fmt.Sprintf(format, a...)
	l.m["level"] = "ERROR"
	l.template.Execute(&buffer, l.m)
	root.queue <- logMessage{msg: buffer.String()}
}

func (l *logger) Fatal(msg string) {
	if l.LogLevel > FATAL {
		return
	}
	buffer := bytes.Buffer{}
	l.m["time"] = time.Now().Format(time.RFC3339Nano)
	l.m["msg"] = msg
	l.m["level"] = "FATAL"
	l.template.Execute(&buffer, l.m)
	root.queue <- logMessage{msg: buffer.String(), fatal: true}
}

func (l *logger) Fatalf(format string, a ...any) {
	if l.LogLevel > FATAL {
		return
	}
	buffer := bytes.Buffer{}
	l.m["time"] = time.Now().Format(time.RFC3339Nano)
	l.m["msg"] = fmt.Sprintf(format, a...)
	l.m["level"] = "FATAL"
	l.template.Execute(&buffer, l.m)
	root.queue <- logMessage{msg: buffer.String(), fatal: true}
}

func (l *logger) Panic(msg string) {
	if l.LogLevel > PANIC {
		return
	}
	buffer := bytes.Buffer{}
	l.m["time"] = time.Now().Format(time.RFC3339Nano)
	l.m["msg"] = msg
	l.m["level"] = "PANIC"
	l.template.Execute(&buffer, l.m)
	root.queue <- logMessage{msg: buffer.String(), panic: true}
}

func (l *logger) Panicf(format string, a ...any) {
	if l.LogLevel > PANIC {
		return
	}
	buffer := bytes.Buffer{}
	l.m["time"] = time.Now().Format(time.RFC3339Nano)
	l.m["msg"] = fmt.Sprintf(format, a...)
	l.m["level"] = "PANIC"
	l.template.Execute(&buffer, l.m)
	root.queue <- logMessage{msg: buffer.String(), panic: true}
}

func (l *logger) SetName(name string) {
	l.name = name
	l.m["logger"] = name
}

func Stop() {
	root.quit <- true
}
