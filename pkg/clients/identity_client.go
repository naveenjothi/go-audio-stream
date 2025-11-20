package clients

import (
	"context"
	"fmt"
	"time"

	pb "go-audio-stream/pkg/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IdentityClient struct {
	client pb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewIdentityClient(host string) (*IdentityClient, error) {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}

	client := pb.NewAuthServiceClient(conn)
	return &IdentityClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *IdentityClient) VerifyToken(token string) (*pb.VerifyTokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return c.client.VerifyToken(ctx, &pb.VerifyTokenRequest{Token: token})
}

func (c *IdentityClient) Close() {
	c.conn.Close()
}
