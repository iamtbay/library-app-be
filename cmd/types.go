package main

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetBookInfo struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Kind        string             `json:"kind" bson:"kind"`
	Destination string             `json:"dst" bson:"dst"`
	Extension   string `json:"extension" bson:"extension"`
}
type NewBookInfo struct {
	Name        string `json:"name" bson:"name"`
	Kind        string `json:"kind" bson:"kind"`
	Destination string `json:"dst" bson:"dst"`
	Extension   string `json:"extension" bson:"extension"`
}

type AuthInfo struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type GetAuthInfo struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password,omitempty" bson:"password"`
	IsAdmin  bool               `json:"isAdmin" bson:"isAdmin"`
}
type NewAuthInfo struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	IsAdmin  bool   `json:"isAdmin" bson:"isAdmin"`
}

type jwtClaims struct {
	UserID   primitive.ObjectID `json:"user_id"`
	Username string             `json:"username"`
	IsAdmin  bool               `json:"isAdmin"`
	jwt.RegisteredClaims
}
