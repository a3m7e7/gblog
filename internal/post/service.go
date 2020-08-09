package post

import (
	"context"
	"github.com/a3m7e7/gblog/pkg/gblog/post"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func New(collection *mongo.Collection) post.PostStorageServer {
	return &postService{
		collection: collection,
	}
}

type postService struct {
	collection *mongo.Collection
}

func (s *postService) GetAllPosts(ctx context.Context,
	request *post.GetAllPostsRequest) (*post.GetAllPostsResponse, error) {

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	results := make([]*post.Post, 0)
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return &post.GetAllPostsResponse{
		Posts: results,
	}, nil
}

func (s *postService) AddPost(ctx context.Context, request *post.AddPostRequest) (*post.AddPostResponse, error) {
	if _, err := s.collection.InsertOne(ctx, request.Post); err != nil {
		return nil, err
	}
	return &post.AddPostResponse{}, nil
}
