package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Layout struct {
	Keys [][]string
	Name string
}

// Stats contains relevant data for a layout
type Stats struct {
	TopSFBS            []SFB
	RowDistribution    []int
	FingerDistribution [8]int
	SFBamount          int
	AlternationAmount  int
	FingerDistance     int
	TrueDistance       float64
	Time               float64
	TextLength         int
	PinkyDistance      int
	OutwardRolls       int
	Redirections       int
	HeatMap            [3][]int
	Score              float64
}

type SFB struct {
	Bigram string
	Count  int
}

// LoadLayouts reads all layouts from the directory into the global Layouts variable.
func LoadLayouts() {
	files, err := ioutil.ReadDir("./layouts/")
	if err != nil {
		panic(err)
	}

	layoutmap := make(map[string]Layout)

	for _, f := range files {
		bytes, err := ioutil.ReadFile("./layouts/" + f.Name())
		if err != nil {
			panic(err)
		}
		s := string(bytes)
		var l Layout
		lines := strings.Split(s, "\n")
		l.Name = lines[0]
		for i := 1; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if line != "" {
				l.Keys = append(l.Keys, strings.Split(line, " "))
			}
		}

		if f.Name() != "optimal" {
			layoutmap[f.Name()] = l
		} else {
			optimal, exists := Layouts["optimal"]
			fmt.Println(exists)
			if !exists {
				layoutmap[f.Name()] = l
			} else {
				layoutmap["optimal"] = optimal
			}
		}

	}

	Layouts = layoutmap
}

// PositionForKey takes in a character as an input, then returns what
// column and row that key is found in on the layout.
func (l *Layout) PositionForKey(char string) (int, int, error) {
	if char == "?" {
		char = "/"
	} else if char == "\"" {
		char = "'"
	} else if char == ":" {
		char = ";"
	}
	for y, row := range l.Keys {
		for x, key := range row {
			if key == char {
				return x, y, nil
			}
		}
	}

	error := fmt.Errorf("There is no key that contains the %s character.", char)
	return 0, 0, error
}
