package js

import "syscall/js"

type Style struct {
	s js.Value
}

func (s *Style) SetWidth(v string) {
	s.setAttribute("width", v)
}

func (s *Style) SetHeight(v string) {
	s.setAttribute("height", v)
}

func (s *Style) SetBackgroundColor(v string) {
	s.setAttribute("background-color", v)
}

func (s *Style) SetBorder(v string) {
	s.setAttribute("border", v)
}

func (s *Style) setAttribute(property string, v interface{}) {
	s.s.Set(property, v)
}
