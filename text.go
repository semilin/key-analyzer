package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// LoadTexts() reads from the ./texts directory and loads all text
// files into the Texts variable.
func LoadTexts() {
	directory, err := os.Open("./texts")
	if err != nil {
		panic("No texts directory was found. Please create one and add at least one sample text file.")
	}
	defer directory.Close()

	files, _ := directory.Readdirnames(0)
	for _, f := range files {
		text, err := os.Open("./texts/" + f)
		if err != nil {
			fmt.Printf("text.go | An error occurred when opening file %s: %s", f, err)
			continue
		}
		content, err := ioutil.ReadAll(text)
		if err != nil {
			//fmt.Printf("text.go | An error occurred when reading file %s: %s", f, err)
			continue
		}

		Texts = append(Texts, string(content))
	}
	FullText = strings.Join(Texts, "  ")
	TextLength = float64(len(strings.ReplaceAll(FullText, " ", "")))
}
