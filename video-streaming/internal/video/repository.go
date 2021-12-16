package video

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MongoClientTimeout = 5
)

type VideoRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func makeContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), MongoClientTimeout*time.Second)
}

func NewVideoRepository(conn string, database string) (*VideoRepository, error) {
	ctx, cancelFunc := makeContext()
	defer cancelFunc()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err != nil {
		return nil, err
	}

	return &VideoRepository{
		client:     client,
		database:   client.Database(database),
		collection: client.Database(database).Collection("videos"),
	}, nil
}

func (r *VideoRepository) GetVideo(id string) (*Video, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := makeContext()
	defer cancelFunc()

	var video Video
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&video); err != nil {
		return nil, err
	}

	return &video, nil
}

func (r *VideoRepository) Disconnect() error {
	ctx, cancelFunc := makeContext()
	defer cancelFunc()

	return r.client.Disconnect(ctx)
}