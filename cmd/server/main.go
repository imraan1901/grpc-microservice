package main

import (
	"log"

	"github.com/imraan1901/grpc-microservice/internal/db"
	"github.com/imraan1901/grpc-microservice/internal/rocket"
)

func run () error {
	// responsible for initializing and 
	// starting our gRPC server
	rocketStore, err := db.New()
	if err != nil {
		return err
	}

	err = rocketStore.Migrate()
	if err != nil {
		log.Println("Failed to run migrations")
		return err
	}

	_ = rocket.New(rocketStore)


	return nil
}

func main() {

	if err := run(); err != nil {
		log.Fatal(err)
	}

}
