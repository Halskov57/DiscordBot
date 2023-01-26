package clients

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	stonkclientv1 "github.com/dietzy1/discordbot/src/proto/stonks/v1"
)

// Return *authserviceclient for example
func NewStonkClient() *stonkclientv1.StonkServiceClient {

	conn, err := grpc.DialContext(
		context.Background(),
		"dns:///0.0.0.0"+os.Getenv("AUTH"),
		//"localhost:9000",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	client := stonkclientv1.NewStonkServiceClient(conn)
	return &client
}
