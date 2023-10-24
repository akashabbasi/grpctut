package main

import (
	"context"
	"flag"
	"log"

	"github.com/akashabbasi/pcbook/pb"
	"github.com/akashabbasi/pcbook/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	serverAddr := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddr)

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	log.Fatal("cannot dial server: ", err)

	laptopClient := pb.NewLaptopServiceClient(conn)
	laptop := sample.NewLaptop()

	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Printf("laptop already exists")
		} else {
			log.Fatal("cannot create laptop", err)
		}
		return
	}

	log.Printf("Laptop created with id: %s", res.Id)
}
