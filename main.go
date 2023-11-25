package main

import (
	"os"

	server "github.com/samricotta/tinyurl/server/http"
)

func main() {
	server.Serve(os.Args[1])
	// import something from server and run it
}
