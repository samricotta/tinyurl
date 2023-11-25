package main

import (
	"fmt"
	"os"

	server "github.com/samricotta/tinyurl/server/http"
)

func main() {
	if err := server.Serve(os.Args[1]); err != nil {
		fmt.Println(err)
	}
}
