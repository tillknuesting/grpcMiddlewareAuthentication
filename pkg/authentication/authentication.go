package authentication

import (
	"context"
	"fmt"
	"time"

	"github.com/form3tech-oss/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// good practice to not use basic types as value in context to avoid collisions
type (
	authenticated     bool
	usernameFromClaim string
)

const (
	Authenticated     authenticated     = false
	UsernameFromClaim usernameFromClaim = ""
)

type Env struct {
	SecretSigningKey []byte
}

// AuthFunc is a middleware (interceptor) that extracts token from header
func (env *Env) AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, err
	}

	username, err := validateToken(token, env.SecretSigningKey)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	// set meta data in context for possible later validation
	newCtx := context.WithValue(ctx, Authenticated, true)
	newCtx = context.WithValue(newCtx, UsernameFromClaim, username)

	return newCtx, nil
}

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username string, signingSecret []byte) (string, error) {
	// Create the Claims
	claims := MyCustomClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
			Issuer:    "AuthFunc",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(signingSecret)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func validateToken(tokenString string, signingSecret []byte) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingSecret, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	}

	return "", fmt.Errorf("could not validate")
}
