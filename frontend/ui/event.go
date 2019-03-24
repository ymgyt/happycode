package ui

import "github.com/ymgyt/happycode/frontend/js"

type KeyboardEvent struct {
	Key       string
	OptionKey bool // or Alt
	CtrlKey   bool
	CmdKey    bool // on Windows windows key
	ShiftKey  bool
	Repeat    bool

	raw *js.Event
}
