package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/imraan1901/grpc-microservice/internal/rocket"
	rkt "github.com/imraan1901/grpc-proto/rocket/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RocketService - defines the interface that the concrete
// implementation has to adhear to
type RocketService interface {
	GetRocketByID(ctx context.Context, id string) (rocket.Rocket, error)
	InsertRocket(ctx context.Context, rkt rocket.Rocket) (rocket.Rocket, error)
	DeleteRocket(ctx context.Context, id string) error
}

// Handler - will handle incoming gRPC requests
type Handler struct {
	RocketService RocketService
	rkt.UnimplementedRocketServiceServer
}

// New - returns a new gRPC handler
func New(rktService RocketService) Handler {
	return Handler{
		RocketService: rktService,
	}
}

func (h Handler) Serve() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Print("could not listen on port 50051")
		return err
	}

	grpcServer := grpc.NewServer()
	rkt.RegisterRocketServiceServer(grpcServer, &h)

	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("failed to serve: %s\n", err)
		return err
	}
	return nil
}

// GetRocket - retrieves a rocket by id and returns the response
func (h Handler) GetRocket(ctx context.Context, req *rkt.GetRocketRequest) (*rkt.GetRocketResponse, error) {
	log.Println("GetRocket gRPC endpoint hit")

	rocket, err := h.RocketService.GetRocketByID(ctx, req.Id)
	if err != nil {
		fmt.Printf("Failed to retrieve rocket by id: %s", err)
		return &rkt.GetRocketResponse{}, nil
	}

	return &rkt.GetRocketResponse{
		Rocket: &rkt.Rocket{
			Id:   rocket.ID,
			Name: rocket.Name,
			Type: rocket.Type,
		},
	}, nil
}

// AddRocket - Adds rocket to the database
func (h Handler) AddRocket(ctx context.Context, req *rkt.AddRocketRequest) (*rkt.AddRocketResponse, error) {
	log.Println("AddRocket gRPC endpoint hit")

	if _, err := uuid.Parse(req.Rocket.Id); err != nil {
		errorStatus := status.Error(codes.InvalidArgument, "uuid is not valid")
		log.Print("given uuid is not valid")
		return &rkt.AddRocketResponse{}, errorStatus
	}

	newRkt, err := h.RocketService.InsertRocket(ctx, rocket.Rocket{
		ID:   req.Rocket.Id,
		Name: req.Rocket.Name,
		Type: req.Rocket.Type,
	})
	if err != nil {
		log.Print("Failed to insert rocket into database")
		return &rkt.AddRocketResponse{}, nil
	}
	return &rkt.AddRocketResponse{
		Rocket: &rkt.Rocket{
			Id:   newRkt.ID,
			Name: newRkt.Name,
			Type: newRkt.Type,
		},
	}, nil
}

// DeleteRocket - Method to delete rocket
func (h Handler) DeleteRocket(ctx context.Context, req *rkt.DeleteRocketRequest) (*rkt.DeleteRocketResponse, error) {
	log.Print("DeleteRocket gRPC endpoint hit")
	err := h.RocketService.DeleteRocket(ctx, req.Id)
	if err != nil {
		return &rkt.DeleteRocketResponse{
			Status: err.Error(),
		}, nil
	}

	return &rkt.DeleteRocketResponse{
		Status: "Successfully deleted rocket",
	}, nil
}
