package app

import (
	"github.com/jopbrown/codesum/pkg/sumer"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	_EVENT_PUSH_MESSAGE  = "pushMessage"
	_EVENT_STREAM_ANSWER = "streamAnswer"
)

func (a *App) PushMessage(msg *sumer.Message) {
	runtime.EventsEmit(a.ctx, _EVENT_PUSH_MESSAGE, msg)
}

func (a *App) StreamAnswer(delta string) {
	runtime.EventsEmit(a.ctx, _EVENT_STREAM_ANSWER, delta)
}
