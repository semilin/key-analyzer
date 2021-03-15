package main

var (
	Texts      []string
	FullText   string
	Data       []Word
	TextLength float64
	Generating bool
	Temp       int
	Layouts    map[string]Layout
	Page       string
)

// placeholders
const (
	NoHand         int = 3
	NoFinger       int = 20
	NoRow          int = 4
	NoDirection    int = 0
	LeftDirection      = -1
	RightDirection     = 1
	LeftHand       int = 0
	RightHand      int = 1
)
