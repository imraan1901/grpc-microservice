package test

import (
	"log"

	"github.com/imraan1901/grpc-proto/rocket/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetClient() rocket.RocketServiceClient {
	var conn *grpc.ClientConn
	// TLS not enabled in test
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("could not connect: %w", err)
	}

	rocketClient := rocket.NewRocketServiceClient(conn)
	return rocketClient

}
