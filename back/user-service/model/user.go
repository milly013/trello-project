package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Email string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	VerificationCode string             `bson:"verificationCode" json:"-"` // za privremeno čuvanje koda
	IsActive bool               `bson:"isActive" json:"isActive"`   // status naloga
}
