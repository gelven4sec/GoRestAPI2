package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var listAttractions []Attraction

func root(w http.ResponseWriter, r *http.Request) {
	// Log
	fmt.Println(time.Now().Format("02-01-2006 15:04:05") + " " + r.Method + " " + r.RequestURI)

	t, _ := os.ReadFile("./index.html")
	fmt.Fprintf(w, string(t))
}

func initServer() {
	fmt.Print("GoRestAPI v1.0\n" +
		"Listining on 0.0.0.0:8000\n")

	http.HandleFunc("/", root)
	http.HandleFunc("/attraction", handleAttractions)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}

func main() {
	listAttractions = initAttractions()
	initServer()
}
