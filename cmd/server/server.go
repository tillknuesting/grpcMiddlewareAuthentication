package main

import (
	"flag"
	"grpcMiddlewareAuth/internal/handlers"
	"grpcMiddlewareAuth/pkg/authentication"
	"log"
	"net"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"google.golang.org/grpc"

	pb "grpcMiddlewareAuth/pkg/proto"
)

func main() {
	var (
		addr   = flag.String("address", ":5050", "address of serverAuthenticated")
		secret = flag.String("secret", "IkWgP3DCsrMJ", "secret to sign tokens")
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
			grpcauth.UnaryServerInterceptor(env.AuthFunc),
		)),
	)

	pb.RegisterGreeterServer(s, &handlers.ServerAuthenticated{})
	pb.RegisterLoginServer(s, &handlers.ServerUnauthenticated{SigningSecret: secretBytes})

	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
