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

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	doBidirectionalStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
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

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
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

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
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

func doBidirectionalStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting a Bidirectional Streaming Call...")
	numbers := [10]int32{1,2,3,4,10,6,7,4,12,5}

	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	waitc := make(chan struct{})

	go func(){
		for _, number := range numbers {
			req := &calculatorpb.FindMaxRequest{
				Number: number,
			}
			fmt.Printf("Sending request: %v", req)
			stream.Send(req)
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v\n", err)
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}