package errcode

type Type int

const (
	OK Type = iota
	FileNotFound
	UnHandled
)

func (t Type) String() string {
	switch t {
	case OK:
		return "OK"
	case FileNotFound:
		return "FILE_NOT_FOUND"
	case UnHandled:
		return "UN_HANDLED"
	}
	return "UNEXPECTED_ERROR"
}
