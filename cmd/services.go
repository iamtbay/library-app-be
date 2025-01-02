package main

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Services struct{}

func initServices() *Services {
	return &Services{}
}

var repo = initRepository()

func (x *Services) getAllBooks() ([]*GetBookInfo, error) {

	return repo.getAllBooks()
}

func (x *Services) addABook(bookInfo NewBookInfo) error {
	return repo.addABook(bookInfo)

}
func (x *Services) getABook(idStr string) (GetBookInfo, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return GetBookInfo{}, nil
	}
	return repo.getABook(id)

}

func (x *Services) deleteABook(idStr string) error {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return err
	}
	return repo.deleteABook(id)
}

//
func (x *Services) login(authInfo AuthInfo) (GetAuthInfo, string, error) {
	userInfo, err := repo.login(authInfo)
	if err != nil {
		return GetAuthInfo{}, "", err
	}
	token, err := createJWT(userInfo)
	return userInfo, token, err
}

func (x *Services) register(authInfo NewAuthInfo) error {
	//check is username unique
	err := repo.isUsernameUnique(authInfo.Username)
	if err != nil {
		return errors.New("username is already in use")
	}
	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(authInfo.Password), 8)
	if err != nil {
		return err
	}
	authInfo.Password = string(hashedPassword)
	//save db
	err = repo.register(authInfo)
	if err != nil {
		return err
	}
	return nil
}
