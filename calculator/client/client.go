package main

import (
	"context"
	"fmt"
	"grpc-greet/calculator/calculatorpb"
	"io"
	"log"
	"time"

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

	// css := calculatorpb.NewPrimeDecompositionServiceClient(cc)
	// doServerStreaming(css)

	ccs := calculatorpb.NewAverageServiceClient(cc)
	doClientStreaming(ccs)
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

func doClientStreaming(c calculatorpb.AverageServiceClient) {
	fmt.Println("Starting a Client Streaming Call...")

	numbers := [4]int32{1,2,3,4}

	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	for _, number := range numbers {
		fmt.Printf("Sending number %v\n", number)
		stream.Send(&calculatorpb.AverageRequest{
			Number: number,
		})
		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving responses: %v", err)
	}
	fmt.Printf("Average: %v\n",res.GetResult())
}