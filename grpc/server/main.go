package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"

	pb "server/grpc/helloworld"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("Received:", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *Server) ServerStreamingSayHello(in *pb.HelloRequest, stream pb.Greeter_ServerStreamingSayHelloServer) error {
	fmt.Println("Received: ", in.GetName())
	for i := 0; i < 10; i++ {
		fmt.Printf("echo message %v\n", in.GetName())
		err := stream.Send(&pb.HelloReply{Message: "hello " + in.GetName()})
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *Server) ClientStreamingSayHello(stream pb.Greeter_ClientStreamingSayHelloServer) error {
	var message string
	for {
		in, err := stream.Recv()
		// 接收结束
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		message += in.GetName()
	}

	fmt.Println("Received: ", message)
	err := stream.SendAndClose(&pb.HelloReply{Message: "hello " + message})
	if err != nil {
		return err
	}
	return nil
}
func (s *Server) BidirectionalStreamingSayHello(stream pb.Greeter_BidirectionalStreamingSayHelloServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Println("bidirectional request received:", in)
		if err := stream.Send(&pb.HelloReply{Message: "bidirectional done"}); err != nil {
			return err
		}
	}
}
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		fmt.Println("failed to listen:", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &Server{})
	fmt.Println("server listening at:", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Println("failed to serve:", err)
	}
}
