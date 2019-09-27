// Package main implements a server for Greeter service.
package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"

	pb "grpcMiddlewareAuth/proto"
)

const (
	port = ":50051"
)

func AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	log.Println(token)
	if token != "testToken" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
	newCtx := context.WithValue(ctx, "basic", token)
	return newCtx, nil
}

type server struct{}

func (s *server) SayHelloAuthenticated(context.Context, *pb.HelloRequest) (*pb.HelloReply, error) {
	panic("implement me")
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(AuthFunc),
		)),
	)
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
