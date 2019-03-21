package js

import "syscall/js"

type Logger struct {
	div js.Value
	doc js.Value
}

func NewLogger() *Logger {
	doc := js.Global().Get("document")
	div := doc.Call("getElementById", "world")
	return &Logger{
		div: div,
		doc: doc,
	}
}

func (l *Logger) Write(p []byte) (n int, err error) {
	node := l.doc.Call("createElement", "div")
	node.Set("innerHTML", string(p))
	l.div.Call("appendChild", node)
	return len(p), nil
}
