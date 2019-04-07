package main

import (
	"./p3/data"
)

func main() {
	//router := p3.NewRouter()
	//if len(os.Args) == 2 {
	//	http.Handle("/", http.FileServer(http.Dir("./frontEnd")))
	//	log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	//} else {
	//	fmt.Println("usage: go run main.go <port number>")
	//}
	data.TestMapToString()
}
