package main

var Qwerty = Layout{
	[3][]string{
		{"q", "w", "e", "r", "t", "y", "u", "i", "o", "p"},
		{"a", "s", "d", "f", "g", "h", "j", "k", "l", ";"},
		{"z", "x", "c", "v", "b", "n", "m", ",", ".", "/"},
	},
	"Qwerty",
}

var Colemak = Layout{
	[3][]string{
		{"q", "w", "f", "p", "g", "j", "l", "u", "y", ";"},
		{"a", "r", "s", "t", "d", "h", "n", "e", "i", "o"},
		{"z", "x", "c", "v", "b", "k", "m", ",", ".", "/"},
	},
	"Colemak",
}

var Colemak_DH = Layout{
	[3][]string{
		{"q", "w", "f", "p", "b", "j", "l", "u", "y", ";"},
		{"a", "r", "s", "t", "g", "m", "n", "e", "i", "o"},
		{"z", "x", "c", "d", "v", "k", "h", ",", ".", "/"},
	},
	"Colemak",
}

var Dvorak = Layout{
	[3][]string{
		{"'", ",", ".", "p", "y", "f", "g", "c", "r", "l"},
		{"a", "o", "e", "u", "i", "d", "h", "t", "n", "s"},
		{";", "q", "j", "k", "x", "b", "m", "w", "v", "z"},
	},
	"Dvorak",
}

var Workman = Layout{
	[3][]string{
		{"q", "d", "r", "w", "b", "j", "f", "u", "p", ";"},
		{"a", "s", "h", "t", "g", "y", "n", "e", "o", "i"},
		{"z", "x", "m", "c", "v", "k", "l", ",", ".", "/"},
	},
	"Workman",
}

var Norman = Layout{
	[3][]string{
		{"q", "w", "d", "f", "k", "j", "u", "r", "l", ";"},
		{"a", "s", "e", "t", "g", "y", "n", "i", "o", "h"},
		{"z", "x", "c", "v", "b", "p", "m", ",", ".", "/"},
	},
	"Workman",
}

var Halmak = Layout{
	[3][]string{
		{"w", "l", "r", "b", "z", ";", "q", "u", "d", "j"},
		{"s", "h", "n", "t", ",", ".", "a", "e", "o", "i"},
		{"f", "m", "v", "c", "/", "g", "p", "x", "k", "y"},
	},
	"Halmak",
}

var Henkaku = Layout{
	[3][]string{
		{"w", "l", "r", "b", "z", ";", "q", "u", "d", "j"},
		{"s", "h", "n", "t", ",", ".", "a", "e", "o", "i"},
		{"f", "m", "v", "c", "/", "g", "p", "x", "k", "y"},
	},
	"Halmak",
}

var BEAKL15 = Layout{
	[3][]string{
		{"j", "y", "o", "-", "k", "g", "c", "m", "n", "z"},
		{"h", "i", "e", "a", "u", "d", "s", "t", "r", "p"},
		{"q", "'", ";", "x", "x", "w", "f", "l", "b", "v"},
	},
	"BEAKL 15",
}

var TNWMLC = Layout{
	[3][]string{
		{"t", "n", "w", "m", "l", "c", "b", "p", "r", "h"},
		{"s", "g", "x", "j", "f", "k", "q", "z", "v", ";"},
		{"e", "a", "d", "i", "o", "y", "u", ",", ".", "/"},
	},
	"TNWMLC",
}

var Optimal = Layout {
	[3][]string{
		{"q", "w", "e", "r", "t", "y", "u", "i", "o", "p"},
		{"a", "s", "d", "f", "g", "h", "j", "k", "l", ";"},
		{"z", "x", "c", "v", "b", "n", "m", ",", ".", "/"},
	},
	"Optimal",
}

func strToLayout(s string) (Layout) {
	switch s {
	case "qwerty":
		return Qwerty
	case "colemak":
		return Colemak
	case "colemak_dh":
		return Colemak_DH
	case "dvorak":
		return Dvorak
	case "workman":
		return Workman
	case "tnwmlc":
		return TNWMLC
	case "halmak":
		return Halmak
	case "beakl":
		return BEAKL15
	case "norman":
		return Norman
	case "optimal":
		return Optimal
	default:
		return Qwerty
	}
}
