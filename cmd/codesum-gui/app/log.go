package app

import (
	"github.com/jopbrown/gobase/log"
	"github.com/wailsapp/wails/v2/pkg/logger"
)

type wLogger struct{}

func NewLogger() logger.Logger { return &wLogger{} }

func (wl *wLogger) Print(message string) {
	log.Print(message)
}

func (wl *wLogger) Trace(message string) {
	log.Debug(message)
}
func (wl *wLogger) Debug(message string) {
	log.Debug(message)
}
func (wl *wLogger) Info(message string) {
	log.Info(message)
}
func (wl *wLogger) Warning(message string) {
	log.Warn(message)
}
func (wl *wLogger) Error(message string) {
	log.Error(message)
}
func (wl *wLogger) Fatal(message string) {
	log.Fatal(message)
}
