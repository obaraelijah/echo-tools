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
	initialized bool
	logger      []*logger
	LogLevel    LogLevel
	lumber      lumberjack.Logger
	format      string
	name        string
	queue       chan logMessage
	quit        chan bool
	template    *template.Template
}

func GetLogLevel() LogLevel {
	if root != nil {
		return root.LogLevel
	}
	return DEBUG
}

func Initialize(config *Config) {
	if root != nil {
		if !root.initialized {
			root.name = "root"
			root.format = "{{.time}} :: {{.level}} : {{.logger}} : {{.msg}}"
			root.queue = make(chan logMessage, config.LogQueue)
			root.quit = make(chan bool)

		}
	} else {
		root = &logger{
			name:   "root",
			format: "{{.time}} :: {{.level}} : {{.logger}} : {{.msg}}",
			queue:  make(chan logMessage, config.LogQueue),
			quit:   make(chan bool),
		}
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

	for _, child := range root.logger {
		child.template = root.template
		child.LogLevel = root.LogLevel
	}

	root.initialized = true

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
	if root == nil {
		root = &logger{}
	}
	l := &logger{
		format:   root.format,
		template: root.template,
		LogLevel: root.LogLevel,
	}
	root.logger = append(root.logger, l)
	l.SetName(name)
	return l
}

func (l *logger) Debug(msg string) {
	if l.LogLevel > DEBUG {
		return
	}
	buffer := bytes.Buffer{}
	m := make(map[string]string)
	m["logger"] = l.name
	for _, line := range strings.Split(msg, "\n") {
		m["time"] = time.Now().Format(time.RFC3339)
		m["msg"] = line
		m["level"] = "DEBUG"
		l.template.Execute(&buffer, m)
		root.queue <- logMessage{msg: buffer.String()}
		buffer.Reset()
	}
}

func (l *logger) Debugf(format string, a ...any) {
	if l.LogLevel > DEBUG {
		return
	}
	buffer := bytes.Buffer{}
	m := make(map[string]string)
	m["logger"] = l.name
	for _, line := range strings.Split(fmt.Sprintf(format, a...), "\n") {
		m["time"] = time.Now().Format(time.RFC3339)
		m["msg"] = line
		m["level"] = "DEBUG"
		l.template.Execute(&buffer, m)
		root.queue <- logMessage{msg: buffer.String()}
		buffer.Reset()
	}
}

func (l *logger) Info(msg string) {
	if l.LogLevel > INFO {
		return
	}
	buffer := bytes.Buffer{}
	m := make(map[string]string)
	m["logger"] = l.name
	for _, line := range strings.Split(msg, "\n") {
		m["time"] = time.Now().Format(time.RFC3339)
		m["msg"] = line
		m["level"] = "INFO "
		l.template.Execute(&buffer, m)
		root.queue <- logMessage{msg: buffer.String()}
		buffer.Reset()
	}

}

func (l *logger) Infof(format string, a ...any) {
	if l.LogLevel > INFO {
		return
	}
	buffer := bytes.Buffer{}
	m := make(map[string]string)
	m["logger"] = l.name
	for _, line := range strings.Split(fmt.Sprintf(format, a...), "\n") {
		m["time"] = time.Now().Format(time.RFC3339)
		m["msg"] = line
		m["level"] = "INFO "
		l.template.Execute(&buffer, m)
		root.queue <- logMessage{msg: buffer.String()}
		buffer.Reset()
	}
}

func (l *logger) Warn(msg string) {
	if l.LogLevel > WARN {
		return
	}
	buffer := bytes.Buffer{}
	m := make(map[string]string)
	m["logger"] = l.name
	for _, line := range strings.Split(msg, "\n") {
		m["time"] = time.Now().Format(time.RFC3339)
		m["msg"] = line
		m["level"] = "WARN "
		l.template.Execute(&buffer, m)
		root.queue <- logMessage{msg: buffer.String()}
		buffer.Reset()
	}
}

func (l *logger) Warnf(format string, a ...any) {
	if l.LogLevel > WARN {
		return
	}
	buffer := bytes.Buffer{}
	m := make(map[string]string)
	m["logger"] = l.name
	for _, line := range strings.Split(fmt.Sprintf(format, a...), "\n") {
		m["time"] = time.Now().Format(time.RFC3339)
		m["msg"] = line
		m["level"] = "WARN "
		l.template.Execute(&buffer, m)
		root.queue <- logMessage{msg: buffer.String()}
		buffer.Reset()
	}
}

func (l *logger) Error(msg string) {
	if l.LogLevel > ERROR {
		return
	}
	buffer := bytes.Buffer{}
	m := make(map[string]string)
	m["logger"] = l.name
	for _, line := range strings.Split(msg, "\n") {
		m["time"] = time.Now().Format(time.RFC3339)
		m["msg"] = line
		m["level"] = "ERROR"
		l.template.Execute(&buffer, m)
		root.queue <- logMessage{msg: buffer.String()}
		buffer.Reset()
	}
}

func (l *logger) Errorf(format string, a ...any) {
	if l.LogLevel > ERROR {
		return
	}
	buffer := bytes.Buffer{}
	m := make(map[string]string)
	m["logger"] = l.name
	for _, line := range strings.Split(fmt.Sprintf(format, a...), "\n") {
		m["time"] = time.Now().Format(time.RFC3339)
		m["msg"] = line
		m["level"] = "ERROR"
		l.template.Execute(&buffer, m)
		root.queue <- logMessage{msg: buffer.String()}
		buffer.Reset()
	}
}

func (l *logger) Fatal(msg string) {
	if l.LogLevel > FATAL {
		return
	}
	buffer := bytes.Buffer{}
	m := make(map[string]string)
	m["logger"] = l.name
	for _, line := range strings.Split(msg, "\n") {
		m["time"] = time.Now().Format(time.RFC3339)
		m["msg"] = line
		m["level"] = "FATAL"
		l.template.Execute(&buffer, m)
		root.queue <- logMessage{msg: buffer.String()}
		buffer.Reset()
	}
}

func (l *logger) Fatalf(format string, a ...any) {
	if l.LogLevel > FATAL {
		return
	}
	buffer := bytes.Buffer{}
	m := make(map[string]string)
	m["logger"] = l.name
	for _, line := range strings.Split(fmt.Sprintf(format, a...), "\n") {
		m["time"] = time.Now().Format(time.RFC3339)
		m["msg"] = line
		m["level"] = "FATAL"
		l.template.Execute(&buffer, m)
		root.queue <- logMessage{msg: buffer.String(), fatal: true}
		buffer.Reset()
	}
}

func (l *logger) Panic(msg string) {
	if l.LogLevel > PANIC {
		return
	}
	buffer := bytes.Buffer{}
	m := make(map[string]string)
	m["logger"] = l.name
	for _, line := range strings.Split(msg, "\n") {
		m["time"] = time.Now().Format(time.RFC3339)
		m["msg"] = line
		m["level"] = "PANIC"
		l.template.Execute(&buffer, m)
		root.queue <- logMessage{msg: buffer.String(), panic: true}
		buffer.Reset()
	}
}

func (l *logger) Panicf(format string, a ...any) {
	if l.LogLevel > PANIC {
		return
	}
	buffer := bytes.Buffer{}
	m := make(map[string]string)
	m["logger"] = l.name
	for _, line := range strings.Split(fmt.Sprintf(format, a...), "\n") {
		m["time"] = time.Now().Format(time.RFC3339)
		m["msg"] = line
		m["level"] = "PANIC"
		l.template.Execute(&buffer, m)
		root.queue <- logMessage{msg: buffer.String(), panic: true}
		buffer.Reset()
	}
}

func (l *logger) SetName(name string) {
	l.name = name
}

func Stop() {
	root.quit <- true
}
