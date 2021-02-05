package main

import (
	"fmt"
	"log"

	"github.com/isongjosiah/grpc/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("I am a client!")
	// creact a client connection to the server
	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error connecting to server : %v", err)
	}
	defer conn.Close()

	c := blogpb.NewBlogServiceClient(conn)
	fmt.Println(c)
}
