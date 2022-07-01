package client

import (
	"log"
	"os"

	pb "work-demo/grpc/protoful/word"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	Address = "localhost:50001"
)

func WordClient() {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHelloWorldServiceClient(conn)
	name := "lin"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Println(r.Message)
}
