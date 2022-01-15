package main

import (
	"context"
	"fmt"
	"grpc-greet/greet/greetpb"
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

	c := greetpb.NewGreetServiceClient(cc)

	// doUnary(c)
	doServerStreaming(c)
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