package main

import (
	"context"
	"fmt"
	"grpc-greet/calculator/calculatorpb"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) CalSum(ctx context.Context, req *calculatorpb.CalSumRequest) (*calculatorpb.CalSumResponse, error) {
	first_number, second_number := req.GetFirstNumber(), req.GetSecondNumber()

	return &calculatorpb.CalSumResponse{
		Result: first_number + second_number,
	}, nil
}

func (*server) PrimeDecomposition(req *calculatorpb.PrimeDecompositionRequest, stream calculatorpb.CalculatorService_PrimeDecompositionServer) error {
	fmt.Printf("Prime decomposition was invoked with %v\n", req)
	var k int32 = 2
	number := req.GetNumber()
	for ; number > 1 ; {
		if number % k == 0 {
			res := &calculatorpb.PrimeDecompositionResponse{
				Result: k,
			}
			if err := stream.Send(res); err != nil {
				log.Fatalf("Error while streaming: %v", err)
				return err
			}
			number = number / k
			time.Sleep(time.Second)
		} else {
			k = k + 1
		}
	}
  
	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	numbers := []int32{}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			var sum int32 = 0
			for _, num := range numbers {
				sum += num
			}
			average := float32(sum)/float32(len(numbers))
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Result: average,
			})
		}
		if err != nil {
			log.Fatalf("Receive stream error: %v", err)
			return err
		}
		numbers = append(numbers, req.GetNumber())
	}
}

func (*server) FindMax(stream calculatorpb.CalculatorService_FindMaxServer) error {
	fmt.Println("Findmax Bidirectional function invoked")

	var max int32 = 0
	c := make(chan int32)
	errc := make(chan error)

	go func(errc chan<- error) {
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				close(c)
				break
			}
	
			if err != nil {
				log.Fatalf("Error while receiving request: %v\n", err)
				errc<-err
			}
	
			c<- req.GetNumber()
		}
	}(errc)

	go func(errc chan<- error) {
		for number := range c {
			if number > max {
				max = number
			}
			msg := &calculatorpb.FindMaxResponse{
				Result: max,
			}
			if err := stream.Send(msg); err != nil {
				log.Fatalf("Error sending message: %v\n", err)
				errc<-err
			}
		}
		errc<-nil
	}(errc)

	return <-errc
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Square Root RPC Called...")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received Negative Number: %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		SquareRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Listening on port 50051")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}