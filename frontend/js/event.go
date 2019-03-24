package js

import (
	"syscall/js"
)

type EventType int

const (
	ClickEvent EventType = iota
	KeyDownEvent
	KeyPressEvent
)

func (et EventType) String() string {
	switch et {
	case ClickEvent:
		return "click"
	case KeyDownEvent:
		return "keydown"
	case KeyPressEvent:
		return "keypress"
	default:
		panic("unexpcted event type")
	}
}

type Event js.Value

func (e Event) Get(key string) js.Value { return js.Value(e).Get(key) }

type EventHandler func(Event)

func (h EventHandler) jsFunc() js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		h(Event(args[0]))
		return nil
	})
}

type KeyboardEvent struct {
	Key       string
	OptionKey bool // or Alt
	CtrlKey   bool
	CmdKey    bool // on Windows windows key
	ShiftKey  bool
	Repeat    bool
}

func (event *Event) ToKeyboard() KeyboardEvent {
	return KeyboardEvent{
		Key:       event.Get("key").String(),
		OptionKey: event.Get("altKey").Bool(),
		CtrlKey:   event.Get("ctrlKey").Bool(),
		CmdKey:    event.Get("metaKey").Bool(),
		ShiftKey:  event.Get("shiftKey").Bool(),
		Repeat:    event.Get("repeat").Bool(),
	}
}
