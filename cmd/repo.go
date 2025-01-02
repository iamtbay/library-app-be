package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct{}

func initRepository() *Repository {
	return &Repository{}
}

func (x *Repository) getAllBooks() ([]*GetBookInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	filter := bson.M{}
	var books []*GetBookInfo
	cursor, err := bookCollection.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return books, errors.New("no books")
		} else {
			return books, err
		}
	}

	for cursor.Next(ctx) {
		var book GetBookInfo
		err := cursor.Decode(&book)
		if err != nil {
			return books, err
		}
		books = append(books, &book)
	}
	return books, nil
}


func (x *Repository) addABook(document NewBookInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := bookCollection.InsertOne(ctx, document)
	if err != nil {
		return err
	}
	return nil
}

func (x *Repository) getABook(id primitive.ObjectID) (GetBookInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var book GetBookInfo
	filter := bson.M{"_id": id}
	err := bookCollection.FindOne(ctx, filter).Decode(&book)
	if err != nil {
		return GetBookInfo{}, err
	}
	return book, nil
}

func (x *Repository) deleteABook(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	filter := bson.M{"_id": id}

	//
	_, err := bookCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (x *Repository) login(authInfo AuthInfo) (GetAuthInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var userInfo GetAuthInfo
	filter := bson.M{"username": authInfo.Username}
	err := authCollection.FindOne(ctx, filter).Decode(&userInfo)
	if err != nil {
		return GetAuthInfo{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(authInfo.Password))
	if err != nil {
		return GetAuthInfo{}, err
	}

	return GetAuthInfo{
		ID:       userInfo.ID,
		Username: userInfo.Username,
		IsAdmin:  userInfo.IsAdmin,
	}, nil
}

func (x *Repository) register(authInfo NewAuthInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err := authCollection.InsertOne(ctx, authInfo)
	if err != nil {
		return err
	}
	return nil
}

func (x *Repository) isUsernameUnique(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	count, err := authCollection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return err
	}
	if count >= 1 {
		return errors.New("username already in use")
	}
	return nil
}
