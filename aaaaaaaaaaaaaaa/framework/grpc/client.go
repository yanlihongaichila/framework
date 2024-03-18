package grpc

import (
	"context"
	"github.com/JobNing/framework/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Client(ctx context.Context, nacosGroup, toService string) (*grpc.ClientConn, error) {
	conn, err := consul.AgentHealthService(ctx, nacosGroup, toService)
	if err != nil {
		return nil, err
	}
	return grpc.Dial(conn, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
