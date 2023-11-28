package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	server "github.com/samricotta/tinyurl/server/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: client get|post <url|tiny_url>")
		os.Exit(1)
	}

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := server.NewTinyURLClient(conn)
	switch os.Args[1] {
	case "get":
		tinyUrl, err := strconv.ParseUint(os.Args[2], 10, 64)
		if err != nil {
			panic(err)
		}
		resp, err := client.Get(context.TODO(), &server.GetRequest{
			TinyUrl: tinyUrl,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("longUrl", resp.LongUrl)

	case "post":
		resp, err := client.Post(context.TODO(), &server.PostRequest{
			LongUrl: os.Args[2],
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("tinyUrl", resp.TinyUrl)

	default:
		fmt.Println("Usage: client get|post <url|tiny_url>")
	}
}
