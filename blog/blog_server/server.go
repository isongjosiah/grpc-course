package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/isongjosiah/hacking/grpc/blog/blogpb"
	_ "github.com/jinzhu/gorm"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	// if we crash the go code, we get the file  name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Hello !")
	// set up a listener
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("error setting up listener : %v", err)
	}

	s := grpc.NewServer()

	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting Server...")
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// wail for control c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//Block until a signal is recieved
	<-ch

	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("End of program")
}
