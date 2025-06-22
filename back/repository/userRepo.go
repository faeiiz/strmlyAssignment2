package repository

import (
	"back/initializers"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       string `bson:"_id,omitempty"`
	Name     string `bson:"name"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type UserRepository interface {
	Create(user User) error
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
}

type userRepo struct {
	collection *mongo.Collection
}

func NewUserRepository() UserRepository {
	return &userRepo{
		collection: initializers.DB.Collection("users"),
	}
}

func (u *userRepo) Create(user User) error {
	_, err := u.collection.InsertOne(context.TODO(), user)
	return err
}

func (u *userRepo) FindByEmail(email string) (*User, error) {
	var user User
	err := u.collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) FindByID(id string) (*User, error) {
	var user User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = u.collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
