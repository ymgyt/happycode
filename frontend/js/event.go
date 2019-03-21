package js

type EventType int

const (
	ClickEvent EventType = iota
	KeyDownEvent
)

func (et EventType) String() string {
	switch et {
	case ClickEvent:
		return "click"
	case KeyDownEvent:
		return "keydown"
	default:
		panic("unexpcted event type")
	}
}
