package main

import (
	"strings"
	"io/ioutil"
	"fmt"
)

var Text string
var TextLen int 

func LoadText() {
	data, err := ioutil.ReadFile("./text.txt")
	Text = string(data)
	Text = strings.TrimSpace(Text)
	Text = strings.ToLower(Text)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}
}

