package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// LoadPage returns the index file
func LoadPage() {
	b, _ := ioutil.ReadFile("./web/index.html")
	Page = string(b)
}

func style() string {
	b, _ := ioutil.ReadFile("./web/style.css")
	return string(b)
}

func response(w http.ResponseWriter, r *http.Request) {
	LoadPage()
	if r.URL.Path[1:] == "favicon.ico" {
		return
	} else if r.URL.Path[1:] == "style.css" {
		fmt.Fprintf(w, style())

	} else if r.URL.Path[1:] == "temp" {
		fmt.Fprintf(w, strconv.Itoa(Temp))

	} else if r.URL.Path[1:] == "start_generation" {
		fmt.Println("Start generation")
		go generateOptimal()
	} else if r.URL.Path[1:] == "layouts" {
		fmt.Println("Getting layouts")
		LoadLayouts()
		j, err := json.Marshal(Layouts)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, string(j))
	} else if len(r.URL.Path[1:]) > 0 {
		fmt.Println(r.URL.Path)
		layout := Layouts[r.URL.Path[1:]]
		stats := layout.DataStats()
		stats.Score = score(stats)
		json, err := json.Marshal(stats)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, string(json))
	} else {
		fmt.Fprintf(w, Page)
	}

}

func webServer() {
	http.HandleFunc("/", response)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
