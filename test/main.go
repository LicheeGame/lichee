package main

import (
	"flag"
	"fmt"
	"net/http"
)

func retappid() (string, error) {
	return "ddd", nil
}

func main() {

	appid := "dddd"
	if true {
		appid, err := retappid()
		fmt.Printf("%s %s", appid, err)
	}
	fmt.Printf("%s", appid)
	port := flag.Int("port", 7456, "server port")
	flag.Parse()
	http.HandleFunc("/ping", helloHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))
}
