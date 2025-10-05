package output

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var (
	currentLevel LogLevel  = LevelInfo
	outputWriter io.Writer = os.Stdout
	errorWriter  io.Writer = os.Stderr
	lineWidth    int       = 80
)

type Outputter struct {
	level LogLevel
	out   io.Writer
}

func NewOutputter(level LogLevel, out io.Writer) *Outputter {
	return &Outputter{
		level: level,
		out:   outputWriter,
	}
}

func ParseLevel(input string) (LogLevel, error) {
	switch true {
	case strings.EqualFold(input, "debug"):
		return LevelDebug, nil
	case strings.EqualFold(input, "info"):
		return LevelInfo, nil
	case strings.EqualFold(input, "warn"):
		return LevelWarn, nil
	case strings.EqualFold(input, "error"):
		return LevelError, nil
	case strings.EqualFold(input, "fatal"):
		return LevelFatal, nil
	}

	return -1, fmt.Errorf("invalid log level: %s", input)
}

func SetLevel(level LogLevel) {
	currentLevel = level
}

func SetOutput(writer io.Writer) {
	outputWriter = writer
}

func SetErrorOutput(writer io.Writer) {
	errorWriter = writer
}

func SetLineWidth(width int) {
	if width > 0 {
		lineWidth = width
	}
}

func (o *Outputter) Debug(format string, args ...interface{}) {
	o.log(LevelDebug, "DEBUG", format, args...)
}

func (o *Outputter) Info(format string, args ...interface{}) {
	o.log(LevelInfo, "INFO", format, args...)
}

func (o *Outputter) Warn(format string, args ...interface{}) {
	o.log(LevelWarn, "WARN", format, args...)
}

func (o *Outputter) Error(format string, args ...interface{}) {
	o.log(LevelError, "ERROR", format, args...)
}

func (o *Outputter) Fatal(format string, args ...interface{}) {
	o.log(LevelFatal, "FATAL", format, args...)
	os.Exit(1)
}

func (o *Outputter) log(level LogLevel, prefix, format string, args ...interface{}) {
	if level < currentLevel {
		return
	}

	message := fmt.Sprintf(format, args...)
	message = o.wrapText(message)

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("%s [%s] %s\n", timestamp, prefix, message)

	if level >= LevelError {
		fmt.Fprint(errorWriter, logEntry)
	} else {
		fmt.Fprint(outputWriter, logEntry)
	}
}

func (o *Outputter) wrapText(text string) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+len(word)+1 > lineWidth {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine += " " + word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

func (o *Outputter) Print(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if !strings.HasSuffix(message, "\n") {
		message = message + "\n"
	}
	fmt.Fprint(o.out, message)
}

func (o *Outputter) Println(format string, args ...interface{}) {
	o.Print(format, args...)
}
