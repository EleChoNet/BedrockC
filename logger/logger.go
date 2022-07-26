package logger

import (
	"fmt"
	"github.com/fatih/color"
	"sync"
	"time"
)

type logger struct {
	lock sync.Mutex
}

var globalLogger logger

func DefaultLogger() *logger {
	return &globalLogger
}
func (l *logger) logging() {
	l.lock.Lock()
	defer l.lock.Unlock()
}
func getNowTime() string {
	return fmt.Sprintf("%d-%d-%d %d:%d:%d", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())

}
func (l *logger) Message(msg string, sender string) {
	l.logging()
	fmt.Printf("[Message][%s]%s:%s\n", getNowTime(), sender, msg)
}
func (l *logger) Warn(msg string, sender string) {
	l.logging()
	color.Set(color.FgYellow)
	fmt.Printf("[Warn][%s]%s:%s\n", getNowTime(), sender, msg)
	color.Unset()
}
func (l *logger) Error(err error, msg string, sender string) {
	l.logging()
	color.Set(color.FgRed)
	fmt.Printf("[Error][%s]%s:%s\n", getNowTime(), sender, msg)
	fmt.Printf("[ErrorStackTrace]:%+v\n", err)
	color.Unset()
}
