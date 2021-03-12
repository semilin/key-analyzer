package main

import (
	"fmt"
)

type Layout struct {
	Keys [3][]string
	Name string
}

type Stats struct {
	TopSFBS            []SFB
	RowDistribution    []int
	FingerDistribution [8]int
	Layout             [3][]string
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
