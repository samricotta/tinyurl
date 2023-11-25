package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: client <url>")
		os.Exit(1)
	}
	url := os.Args[1]

	resp, err := http.Post("http://localhost:8080", "text/plain", strings.NewReader(url))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("tinyUrl", string(body))
}
