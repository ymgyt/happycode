package ui

import (
	"github.com/ymgyt/happycode/core/config"
	"github.com/ymgyt/happycode/frontend/js"
	"github.com/ymgyt/happycode/frontend/log"
)

type UI struct {
	Document *js.Element
	World    *js.Element
}

func (ui *UI) ApplyTheme(theme *config.Theme) {
	log.V(10).Debug("apply theme")
	s := ui.World.Style()
	s.SetWidth("100%")
	s.SetHeight("100%")
	s.SetBackgroundColor(theme.BackgroundColor)
}
