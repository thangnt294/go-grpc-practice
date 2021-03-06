package main

import (
	"context"
	"fmt"
	"grpc-greet/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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
	// doClientStreaming(c)
	// doBidirectionalStreaming(c)
	doUnaryWithDeadline(c, 5 * time.Second) // should complete
	doUnaryWithDeadline(c, 1 * time.Second) // should timeout
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

func doBidirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting a Bidirectional Streaming Call...")
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	request := []*greetpb.GreetEveryoneRequest{
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

	waitc := make(chan struct{})
	go func() {
		for _, req := range request {
			fmt.Printf("Sending message: %v\n", err)
			stream.Send(req)
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()

	go func(){
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

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting a Unary With Deadline Call...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Thomas",
			LastName: "Nguyen",
		},
	})

	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			if resErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("Unexpected error: %v", err)
			}
		} else {
			log.Fatalf("Could not greet with deadline: %v", err)
		}
	} else {
		fmt.Println(res.GetResult());
	}

}