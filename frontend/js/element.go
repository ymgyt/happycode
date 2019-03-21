package js

import (
	"syscall/js"
)

type Element struct {
	e js.Value
}

func NewElement(element js.Value) *Element {
	assertType(element, TypeObject)
	return &Element{e: element}
}

func (e *Element) GetElementByID(id string) (*Element, error) {
	element := e.e.Call("getElementById", id)
	if isNull(element) {
		return nil, ErrElementNotFound
	}
	return NewElement(element), nil
}

func (e *Element) AddEventListener(et EventType, cb js.Func) {
	e.e.Call("addEventListener", et.String(), cb)
}

func (e *Element) Style() *Style {
	return &Style{s: e.e.Get("style")}
}

