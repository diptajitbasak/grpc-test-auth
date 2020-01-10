package main

import (
	"context"
	"net"
	"fmt"
	"errors"
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
	Id string `json:"id"`
	jwt.StandardClaims
}

func SaveToDB(token string, machId string) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	client.Set(machId, token, 0)
	if err != nil {
		log.Fatal(err)
	}
	// Output: PONG <nil>
}

func CheckFromDB(token string, machId string) bool {
	var ret bool;
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	val, err := client.Get(machId).Result()
	if err != nil {
		fmt.Println("machineId not found")
		ret = false
	}
	if val == token {
		fmt.Println("databse respoonse verified");
		ret = true;
	} else {
		ret = false;
	}
	return ret;
}

func (s *server) GenToken (ctx context.Context, request *protos.Request) (*protos.Response, error) {
	machineID := request.GetMachineId()
	claims := &Claims{
		Id: machineID,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	SaveToDB(tokenString, machineID)
	
	return &protos.Response{ Token: tokenString }, err
}

func (s *server) VerifyToken (ctx context.Context, request *protos.VerifyRequest) (*protos.VerifyResponse, error) {
	machineID := request.GetMachineId()
	Token := request.GetToken()
	resp := CheckFromDB(Token, machineID);
	fmt.Println("server response", resp)
	var msg string;
	var err error;
	if resp {
		msg = "Verified!!!"
	} else {
		msg = "Token Invalid!!!"
		err = errors.New("Invalid Token")
	}
	fmt.Println(msg)
	return &protos.VerifyResponse{ Msg: msg }, err
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