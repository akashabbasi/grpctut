package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/akashabbasi/pcbook/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LaptopServer is the server that provides laptop service
type LaptopServer struct {
	Store LaptopStore
	pb.UnimplementedLaptopServiceServer
}

// NewLaptopServer returns a new laptop server
func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{
		Store: store,
	}
}

func (server *LaptopServer) CreateLaptop(
	ctx context.Context,
	req *pb.CreateLaptopRequest,
) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("receive a create-laptop request with id: %s", laptop.Id)

	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"laptop id is not a valid UUID: %v", err,
			)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				"cannot generate a laptop id: %v",
				err,
			)
		}

		laptop.Id = id.String()
	}

	time.Sleep(6 * time.Second)
	if ctx.Err() == context.Canceled {
		log.Printf("request is canceled")
		return nil, status.Error(codes.Canceled, "request is canceled")
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	// save the laptop in database
	err := server.Store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to the store: %v", err)
	}

	log.Printf("saved laptop with id: %s", laptop.Id)
	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}

	return res, nil
}
