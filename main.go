package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}
	fmt.Fprint(w, "hello")
}
func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
	}
	fmt.Fprint(w, "POST request succefull")
    name := r.FormValue("name")
    address := r.FormValue("address")
    fmt.Fprintf(w, "Name = %s\n", name)
    fmt.Fprintf(w, "Address = %s\n", address)
}
func testHandler(w http.ResponseWriter, r *http.Request) {
if r.URL.Path != "/test" {
        http.Error(w, "404 not found", http.StatusNotFound)
        return
    }
    if r.Method != "GET" {
        http.Error(w, "Method is not supported", http.StatusNotFound)
        return
    }
    http.ServeFile(w, r, "./static/test.html")
}
func main() {
	var fileServer = http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)
    http.HandleFunc("/test", testHandler)
	fmt.Print("The server has started at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
