package handlers

import (
	"context"
	"grpcMiddlewareAuth/authentication"
	pb "grpcMiddlewareAuth/proto"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//implements methods that are intercepted by AuthFunc
type ServerAuthenticated struct{}

//implements methods that are intercepted by AuthFuncOverride
type ServerUnauthenticated struct {
	SigningSecret []byte
}

//SayHelloAuthenticated only can be called when a context is authorized
func (s *ServerAuthenticated) SayHelloAuthenticated(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	//this check can be seen as optional but its recommended to double check that the context was authenticated
	authCrossCheck, ok := ctx.Value(authentication.Authenticated).(bool)
	if !ok && !authCrossCheck {
		return nil, status.Errorf(codes.Unauthenticated, "context not authenticated")
	}

	log.Printf("Received: %v", in.GetName())

	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

//GetToken generates a pseudo token and encodes it with base64
func (s *ServerUnauthenticated) GetToken(ctx context.Context, in *pb.GetTokenRequest) (*pb.GetTokenResponse, error) {
	tokenSigned, err := authentication.GenerateToken(in.User, s.SigningSecret)

	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "could not generate token")
	}

	return &pb.GetTokenResponse{Token: tokenSigned}, nil
}

//AuthFuncOverride is called instead of AuthFunc for methods on ServerUnauthenticated struct
func (s *ServerUnauthenticated) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	log.Println("client is calling method:", fullMethodName)

	//this is optional as it should no be possible to get here anyway
	//if fullMethodName != "/helloworld.Login/GetToken" {
	//	return nil, status.Errorf(codes.Unauthenticated, "no token")
	//}

	return ctx, nil
}
