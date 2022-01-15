package main

import (
	"context"
	"fmt"
	"grpc-greet/calculator/calculatorpb"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main () {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	// c := calculatorpb.NewCalSumServiceClient(cc)
	// doUnary(c)

	css := calculatorpb.NewPrimeDecompositionServiceClient(cc)
	doServerStreaming(css)
}

func doUnary(c calculatorpb.CalSumServiceClient) {
	fmt.Println("Starting a Unary Call...")
	res, err := c.CalSum(context.Background() , &calculatorpb.CalSumRequest{
		FirstNumber: 10,
		SecondNumber: 19,
	})

	if err != nil {
		log.Fatalf("Could not calculate sum: %v", err)
	}

	fmt.Println(res);
}

func doServerStreaming(c calculatorpb.PrimeDecompositionServiceClient) {
	fmt.Println("Starting a Server Streaming Call...")
	resStream, err := c.PrimeDecomposition(context.Background(), &calculatorpb.PrimeDecompositionRequest{
		Number: 120,
	})
	if err != nil {
		log.Fatalf("Error while calling server streaming: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		fmt.Println(msg.GetResult())
	}
}