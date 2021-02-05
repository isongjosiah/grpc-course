package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/isongjosiah/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	//dail the server with grpc.Dail to get a client connection
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't establish a connection to the server: %v", err)
	}

	// create a new serviceclient with the client connection obtained
	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnary(c)
	// doServerStream(c)
	// doClientStream(c)
	// doBiDiStream(c)
	doErrorUnary(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	res, err := c.Calculate(context.Background(), &calculatorpb.CalculatorRequest{
		Numbers: &calculatorpb.Numbers{
			NumberOne: 3,
			NumberTwo: 12,
		},
	})

	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}

	log.Printf("Response: %v", res.Result)
}

func doServerStream(c calculatorpb.CalculatorServiceClient) {
	resStream, err := c.Decompose(context.Background(), &calculatorpb.DecomposeRequest{
		PrimeNumber: 210,
	})
	if err != nil {
		log.Fatalf("error while calling Decompose: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error when reading stream: %v", err)
		}
		log.Printf("Response from Decompose: %v", msg.GetResult())
	}
}

func doClientStream(c calculatorpb.CalculatorServiceClient) {
	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 5,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 10,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 15,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 25,
		},
	}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error when calling compute average: %v", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending Request : %v\n", req.GetNumber())
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error when collecting response : %v", err)
	}

	log.Printf("Response from server: %v", res)
}

func doBiDiStream(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error reading from stream : %v", err)
	}

	waitc := make(chan struct{})

	number := []int32{3, 5, 5, 4, 6, 3, 9, 10, 200, 900, 100}
	// sending stuff to the server
	go func() {
		for _, n := range number {
			fmt.Printf("Sending Request: %v\n", n)
			err := stream.Send(&calculatorpb.FindMaximumRequest{
				Number: n,
			})
			if err != nil {
				log.Fatalf("error streaming to server :%v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("error closing stream : %v", err)
		}
	}()

	// reading from the server
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				break
			}
			if err != nil {
				log.Fatalf("error when streaming from server : %v", err)
			}
			log.Printf("Maximum value is: %v\n", res.GetMaximum())
		}
	}()

	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a squarroot unary RPC")
	number := int32(-3)
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: number})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error form gRPC (USER ERROR)
			fmt.Println(respErr.Message())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("we sent a negative number")
				return
			}
		} else {
			log.Fatalf("Error calling square Root: %v\n", err)
			return
		}
	}

	fmt.Printf("Result of square root of %v is %v\n", number, res.GetNumberRoot())
}
