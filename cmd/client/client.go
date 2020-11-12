package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "grpcMiddlewareAuth/pkg/proto"
)

func getToken(conn *grpc.ClientConn, user string, pass string) (string, error) {
	c := pb.NewLoginClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	r, err := c.GetToken(ctx, &pb.GetTokenRequest{
		User:     user,
		Password: pass,
	})

	if err != nil {
		return "", err
	}

	return r.Token, nil
}

func greetWithToken(conn *grpc.ClientConn, token string, name string) (string, error) {
	c := pb.NewGreeterClient(conn)

	// Create metadata and context.
	md := metadata.Pairs("authorization", "bearer "+token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Make RPC using the context with the metadata.
	var header, trailer metadata.MD
	r, err := c.SayHelloAuthenticated(ctx, &pb.HelloRequest{Name: name}, grpc.Header(&header), grpc.Trailer(&trailer))

	if err != nil {
		return "", err
	}

	return r.Message, nil
}

func main() {
	var (
		user = flag.String("user", "test", "username")
		pass = flag.String("password", "test", "password of user")
		name = flag.String("name", "Alice", "name for greeting")
		addr = flag.String("address", ":5050", "address of server")
	)

	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	// getting the token from server
	token, err := getToken(conn, *user, *pass)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Token: %v", token)

	// greeting the server with auth token
	m, err := greetWithToken(conn, token, *name)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("message:", m)
}
