// Package main implements a serverAuthenticated for Greeter service.
package main

import (
	"flag"
	"grpcMiddlewareAuth/authentication"
	"grpcMiddlewareAuth/handlers"
	"log"
	"net"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"google.golang.org/grpc"

	pb "grpcMiddlewareAuth/proto"
)

func main() {
	var (
		addr   = flag.String("address", ":50051", "address of serverAuthenticated")
		secret = flag.String("secret", "randomString", "secret to sign Tokens")
	)

	flag.Parse()

	lis, err := net.Listen("tcp", *addr)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	secretBytes := []byte(*secret)
	env := authentication.Env{SecretSigningKey: secretBytes}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(env.AuthFunc),
		)),
	)

	pb.RegisterGreeterServer(s, &handlers.ServerAuthenticated{})
	pb.RegisterLoginServer(s, &handlers.ServerUnauthenticated{SigningSecret: secretBytes})

	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
