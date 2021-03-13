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

type Stats struct {
	TopSFBS            []SFB
	RowDistribution    []int
	FingerDistribution [8]int
	Layout             [][]string
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
		for i:=1;i<len(lines);i++ {
			l.Keys = append(l.Keys, strings.Split(lines[i], " "))
		}

		layoutmap[f.Name()] = l
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
