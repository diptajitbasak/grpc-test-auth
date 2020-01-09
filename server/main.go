package main

import (
	"context"
	"net"
	"fmt"
	"log"
	"gitlab.com/grpc-test-auth/protos"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {}
var jwtKey = []byte("qwerty")
type Claims struct {
	id string `json:"id"`
	jwt.StandardClaims
}

func ExampleNewClient(token string) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	client.Do("Token", token)
	if err != nil {
		log.Fatal(err)
	}

	// Output: PONG <nil>
}

func main() {
	listener, err := net.Listen("tcp", ":4040");
	if err != nil {
		panic(err)
	}
	
	srv := grpc.NewServer()
	protos.RegisterGenJWTTokenServer(srv, &server{})
	reflection.Register(srv)
	if e := srv.Serve(listener); e != nil {
		panic(err)
	}
}

func (s *server) GenToken (ctx context.Context, request *protos.Request) (*protos.Response, error) {
	machineID := request.GetMachineId()
	claims := &Claims{
		id: machineID,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	ExampleNewClient(tokenString)
	fmt.Println("bibjibiuj",claims)
	return &protos.Response{ Token: tokenString }, err
}