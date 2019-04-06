package main

import (
	"./p3"
	"./p3/data"
	"log"
	"net/http"
	"os"
)

func main() {
	data.TestPeerListRebalance()
	router := p3.NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":6686", router))
	}
}
