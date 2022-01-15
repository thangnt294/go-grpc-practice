package main

import (
	"context"
	"fmt"
	"grpc-greet/greet/greetpb"
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

	c := greetpb.NewGreetServiceClient(cc)

	// doUnary(c)
	// doServerStreaming(c)
	doClientStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting a Unary Call...")
	res, err := c.Greet(context.Background() , &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Thomas",
			LastName: "Nguyen",
		},
	})

	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}

	fmt.Println(res);
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting a Server Streaming Call...")
	resStream, err :=	c.GreetManyTimes(context.Background(), &greetpb.GreetManyTimesRequest{
	Greeting: &greetpb.Greeting{
			FirstName: "Thomas",
			LastName: "Nguyen",
		},
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

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting a Client Streaming Call...")

	request := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Thomas",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Lily",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "George",
			},
		},
	}

	reqStream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling long greet %v", err)
	}

	for _, req := range request {
		fmt.Printf("Sending request: %v\n", req)
		reqStream.Send(req)
		time.Sleep(time.Second)
	}
	res, err:= reqStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving responses: %v", err)
	}
	fmt.Printf("Response: %v\n", res)
}