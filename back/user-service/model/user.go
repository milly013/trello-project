package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username         string             `bson:"username" json:"username"`
	Email            string             `bson:"email" json:"email"`
	Password         string             `bson:"password" json:"password"`
	VerificationCode string             `bson:"verificationCode" json:"-"`
	IsActive         bool               `bson:"isActive" json:"isActive"`

	Role               string             `bson:"role" json:"role"`
    ResetToken         string             `bson:"resetToken,omitempty" json:"resetToken,omitempty"`
    ResetTokenExpiresAt time.Time          `bson:"resetTokenExpiresAt,omitempty" json:"resetTokenExpiresAt,omitempty"`

	

}
