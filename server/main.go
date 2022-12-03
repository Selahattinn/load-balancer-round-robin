package main

import (
	"flag"
	"net/http"
)

var (
	addr = flag.String("h", "localhost:8000", "The host and port of server")
)

func main() {
	flag.Parse()

	if addr == nil {
		panic("addr can  not be nil")
	}

	http.HandleFunc("/", hello)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		panic(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello\n"))
}
