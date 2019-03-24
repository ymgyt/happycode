package js

import (
	"fmt"
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

func (e *Element) AddEventHandler(et EventType, h EventHandler) {
	e.e.Call("addEventListener", et.String(), h.jsFunc())
}

func (e *Element) Style() *Style {
	return &Style{s: e.e.Get("style")}
}

func (e *Element) AppendText(s string) {
	org := e.e.Get("textContent").String()
	fmt.Println("org", org)

	innerText := org + s
	e.e.Set("textContent", innerText)
}

func (e *Element) CreateElement(tag string) *Element {
	created := e.e.Call("createElement", tag)
	return &Element{e: created}
}

func (e *Element) AppendChild(child *Element) {
	assertType(child.e, TypeObject)
	e.e.Call("appendChild", child.e)
}
