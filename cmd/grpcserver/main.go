package main

import (
	"context"
	"log"
	"net"

	"github.com/mzzz-zzm/go-tdd-practice/adapters/grpcserver"
	"github.com/mzzz-zzm/go-tdd-practice/domain/interactions"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	s := grpc.NewServer()
	grpcserver.RegisterGreeterServer(s, &GreetServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type GreetServer struct {
	grpcserver.UnimplementedGreeterServer
}

func (g *GreetServer) Greet(ctx context.Context, request *grpcserver.GreetRequest) (*grpcserver.GreetReply, error) {
	return &grpcserver.GreetReply{
		Message: interactions.Greet(request.Name),
	}, nil
}

func (g *GreetServer) Curse(ctx context.Context, request *grpcserver.GreetRequest) (*grpcserver.GreetReply, error) {
	return &grpcserver.GreetReply{
		Message: interactions.Curse(request.Name),
	}, nil
}
