package repository

import (
	"back/initializers"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Video struct {
	ID          string    `bson:"_id,omitempty"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	URL         string    `bson:"url"`
	UploaderID  string    `bson:"uploader_id"`
	UploadDate  time.Time `bson:"upload_date"`
}

type VideoRepository interface {
	Create(video Video) error
	GetAll() ([]Video, error)
}

type videoRepo struct {
	collection *mongo.Collection
}

func NewVideoRepository() VideoRepository {
	return &videoRepo{
		collection: initializers.DB.Collection("videos"),
	}
}

func (v *videoRepo) Create(video Video) error {
	_, err := v.collection.InsertOne(context.TODO(), video)
	return err
}
func (v *videoRepo) GetAll() ([]Video, error) {
	var videos []Video

	cursor, err := v.collection.Find(context.TODO(), bson.M{}, options.Find().SetSort(bson.M{"upload_date": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var video Video
		if err := cursor.Decode(&video); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}
