package main

import (
	"./p3"
	"./p3/data"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	data.TestPeerListRebalance()
	router := p3.NewRouter()
	if len(os.Args) == 2 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		fmt.Println("usage: go run main.go <port number>")
	}
}
