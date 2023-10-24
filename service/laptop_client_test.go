package service_test

import (
	"context"
	"net"
	"testing"

	"github.com/akashabbasi/pcbook/pb"
	"github.com/akashabbasi/pcbook/sample"
	"github.com/akashabbasi/pcbook/serializer"
	"github.com/akashabbasi/pcbook/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClientCreateLaptop(t *testing.T) {
	laptopServer, serverAddr := startTestLaptopServer(t)
	laptopClient := newTestLaptopClient(t, serverAddr)

	laptop := sample.NewLaptop()
	expectedId := laptop.Id
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedId, res.Id)

	// check that laptop is stored to the database
	other, err := laptopServer.Store.Find(res.Id)
	require.NoError(t, err)
	require.NotNil(t, other)

	// check saved laptop is saame as we sent
	requireSameLaptop(t, laptop, other)
}

func startTestLaptopServer(
	t *testing.T,
) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(
		service.NewInMemoryLaptopStore(),
	)

	grpcServer := grpc.NewServer()

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go grpcServer.Serve(listener) // non blocking call

	return laptopServer, listener.Addr().String()
}

func newTestLaptopClient(
	t *testing.T,
	serverAddr string,
) pb.LaptopServiceClient {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)
}

func requireSameLaptop(t *testing.T, laptop1, laptop2 *pb.Laptop) {
	json1, err := serializer.ProtobufToJSON(laptop1)
	require.NoError(t, err)

	json2, err := serializer.ProtobufToJSON(laptop2)
	require.NoError(t, err)

	require.Equal(t, json1, json2)
}
