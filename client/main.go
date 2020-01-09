package main

import (
	"fmt"
	"gitlab.com/grpc-test-auth/protos"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := protos.NewGenJWTTokenClient(conn)

	g := gin.Default()
	g.GET("/getToken/:machineId", func(ctx *gin.Context) {
		machineId := ctx.Param("machineId")
		// if err!= nil {
		// 	ctx.JSON(http.statusBadRequest, gin.H{"error": "Invalid Parameter machineId"})
		// }
		fmt.Println(machineId)
		req := &protos.Request{MachineId: machineId}
		if response, err := client.GenToken(ctx, req); err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"Token": fmt.Sprint(response.Token),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})
	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}