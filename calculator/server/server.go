package main

import (
	"context"
	"fmt"
	"grpc-greet/calculator/calculatorpb"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) CalSum(ctx context.Context, req *calculatorpb.CalSumRequest) (*calculatorpb.CalSumResponse, error) {
	first_number, second_number := req.GetFirstNumber(), req.GetSecondNumber()

	return &calculatorpb.CalSumResponse{
		Result: first_number + second_number,
	}, nil
}

func (*server) PrimeDecomposition(req *calculatorpb.PrimeDecompositionRequest, stream calculatorpb.PrimeDecompositionService_PrimeDecompositionServer) error {
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

func main() {
	fmt.Println("Listening on port 50051")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalSumServiceServer(s, &server{})
	calculatorpb.RegisterPrimeDecompositionServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}