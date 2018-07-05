package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ogiekako/laundry/ticker"
)

func drawGraph(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		filename := strings.TrimPrefix(r.URL.Path, "/")
		if filename == "" {
			filename = "graph.html"
		}
		body, err := ioutil.ReadFile("static/" + filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Write(body)
	}
}

func handleData(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ticker.Retrieve()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		marshaled, err := json.Marshal(*data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(marshaled))
		return
	} else if r.Method == "POST" {
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ss := strings.Split(string(bs), ",")
		shake := 0
		if !strings.Contains("No", ss[1]) {
			shake, err = strconv.Atoi(ss[1])
			if err != nil {
				shake = 0
			}
		}
		tap := 0
		if !strings.Contains("No", ss[2]) {
			tap, err = strconv.Atoi(ss[2])
			if err != nil {
				tap = 0
			}
		}
		ticker.Add(shake, tap)
		return
	}
}

func main() {
	http.HandleFunc("/", drawGraph)
	http.HandleFunc("/data", handleData)

	ticker.Start()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
