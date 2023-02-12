package clients

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	stonkclientv1 "github.com/dietzy1/discordbot/src/proto/stonk/v1"
)

// Return *authserviceclient for example
func NewStonkClient() *stonkclientv1.StonkServiceClient {

	conn, err := grpc.DialContext(
		context.Background(),
		"dns:///0.0.0.0"+":8000",
		//"localhost:8000",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	log.Println("Connected to stonks service")

	client := stonkclientv1.NewStonkServiceClient(conn)

	return &client
}
