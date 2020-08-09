package main

import (
	"context"
	"github.com/a3m7e7/gblog/internal/post"
	postPb "github.com/a3m7e7/gblog/pkg/gblog/post"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	ctx := context.Background()

	mongoClient := getMongodbClient(ctx)
	postCollection := mongoClient.Database("gblog").Collection("post")

	postService := post.New(postCollection)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	postPb.RegisterPostStorageServer(grpcServer, postService)

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc")
	}
}

func getMongodbClient(ctx context.Context) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("could not create mongo client")
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("could not connect mongo client")
	}
	return client
}
