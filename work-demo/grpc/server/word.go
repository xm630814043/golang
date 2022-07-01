package server

import (
    "google.golang.org/grpc"
    "log"
    "net"
    pb "work-demo/grpc/protoful/word"
    "golang.org/x/net/context"
)

const (
    PORT = ":50001"
)

type Server struct {}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
    log.Println("request: ", in.Name)
    return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func WordServer() {
    lis, err := net.Listen("tcp", PORT)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterHelloWorldServiceServer(s, &Server{})
    log.Println("rpc服务已经开启")
    s.Serve(lis)
}