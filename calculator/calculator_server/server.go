package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/isongjosiah/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Calculate(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	one := req.GetNumbers().GetNumberOne()
	two := req.GetNumbers().GetNumberTwo()
	sum := one + two
	res := &calculatorpb.CalculatorResponse{
		Result: sum,
	}
	return res, nil
}

func (*server) Decompose(req *calculatorpb.DecomposeRequest, stream calculatorpb.CalculatorService_DecomposeServer) error {
	k := int32(2)
	n := req.GetPrimeNumber()

	for n > 1 {
		if n%k == 0 {
			stream.Send(&calculatorpb.DecomposeResponse{
				Result: k,
			})
			n = n / k
		} else {
			k = k + 1
			fmt.Printf("The divisor has increased to %v\n", k)
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	sum := float64(0)
	n := float64(0)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			res := sum / n
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: res,
			})
		}
		if err != nil {
			log.Fatalf("Error while streaming client request : %v", err)
		}
		sum += float64(msg.GetNumber())
		n += 1.0
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("FindMaximum has been invoked")
	var number []int32
	var max int32

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			stream.Send(&calculatorpb.FindMaximumResponse{
				Maximum: max,
			})
			break
		}
		if err != nil {
			log.Fatalf("error streaming from client : %v", err)
			return err
		}
		number = append(number, req.GetNumber())
		if max == max2(number) {
			continue
		}
		max = max2(number)

		err = stream.Send(&calculatorpb.FindMaximumResponse{
			Maximum: max,
		})

		if err != nil {
			log.Fatalf("error sending stream to client : %v", err)
		}
	}
	return nil
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Recieved a negative number")
	}

	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func max2(array []int32) int32 {
	max := array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
	}
	return max
}

func main() {
	// start a new listener
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// start a new grpc server
	s := grpc.NewServer()

	// register that grpc server with the service server
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	//register reflection service on gRPC server
	reflection.Register(s)

	//serve the server. Hehe
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
