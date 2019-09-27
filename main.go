// Package main implements a server for Greeter service.
package main

import (
	"context"
	"encoding/base64"
	"flag"
	"log"
	"net"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"

	pb "grpcMiddlewareAuth/proto"
)

var (
	addr = flag.String("address", ":50051", "address of server")
)

//AuthFunc is a middleware (interceptor) that extracts token from header
func AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	auth := "test" + ":" + "test"
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))
	if token != authEncoded {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
	//note that its not good practice to use basic types as value in context but
	//we are going to do this anyway bc we do not care about collision for demo
	newCtx := context.WithValue(ctx, "authenticated", true)
	return newCtx, nil
}

//implements methods that are intercepted by AuthFunc
type server struct{}

//implements methods that are intercepted by AuthFuncOverride
type authServer struct{}

//SayHelloAuthenticated only can be called when a context is authorized
func (s *server) SayHelloAuthenticated(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

//GetToken generates a pseudo token and encodes it with base64
func (s *authServer) GetToken(ctx context.Context, in *pb.GetTokenRequest) (*pb.GetTokenResponse, error) {
	auth := in.User + ":" + in.Password
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))
	return &pb.GetTokenResponse{Token: authEncoded}, nil
}

//AuthFuncOverride is called instead of AuthFunc for methods on authServer struct
func (s *authServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	log.Println("client is calling method:", fullMethodName)

	//this is optional as it should no be possible to get anyway
	if fullMethodName != "/helloworld.Login/GetToken" {
		return nil, status.Errorf(codes.Unauthenticated, "no auth token used")
	}
	return ctx, nil
}

func main() {
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(AuthFunc),
		)),
	)
	pb.RegisterGreeterServer(s, &server{})
	pb.RegisterLoginServer(s, &authServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
