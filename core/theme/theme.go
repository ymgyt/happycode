package theme

type Theme struct {
	Meta   Meta
	Window Window
}

type Meta struct {
	Name string
}

type Window struct {
	BackgroundColor string
}
