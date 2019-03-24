package ui

import (
	"github.com/ymgyt/happycode/core/config"
	"github.com/ymgyt/happycode/frontend/js"
)

type UI struct {
	Document       *js.Element
	World          *js.Element
	KeyboardEvents chan js.KeyboardEvent

	windows map[string]*Window
	theme   *config.Theme
}

func (ui *UI) Init(cfg *config.Config) {
	ui.theme = cfg.Theme
	ui.InitWindow()
	ui.InitEventListener()
}

func (ui *UI) InitWindow() {
	if ui.windows == nil {
		ui.windows = make(map[string]*Window)
	}
	s := ui.World.Style()
	s.SetWidth("100%")
	s.SetHeight("100%")

	w1 := ui.newMainWindow()
	ui.World.AppendChild(w1.Element)

	debugw := ui.newDebugWindow()
	ui.World.AppendChild(debugw.Element)
}

func (ui *UI) newMainWindow() *Window {
	w := ui.Document.CreateElement("div")
	s := w.Style()
	s.SetWidth("100%")
	s.SetHeight("70%")
	s.SetBackgroundColor(ui.theme.BackgroundColor)
	s.SetBorder("dashed red")

	return &Window{
		Name:    "window-a",
		Element: w,
	}
}

func (ui *UI) newDebugWindow() *Window {
	w := ui.Document.CreateElement("div")
	s := w.Style()
	s.SetWidth("100%")
	s.SetHeight("30%")
	s.SetBackgroundColor(ui.theme.BackgroundColor)
	s.SetBorder("dashed red")

	return &Window{
		Name:    "debug",
		Element: w,
	}
}

func (ui *UI) InitEventListener() {
	ui.Document.AddEventHandler(js.KeyPressEvent, func(event js.Event) {
		ui.KeyboardEvents <- event.ToKeyboard()
	})
}
