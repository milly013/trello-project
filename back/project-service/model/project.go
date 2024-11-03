// model/project.go
package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	EndDate    time.Time          `bson:"endDate" json:"endDate"`
	MinMembers int                `bson:"minMembers" json:"minMembers"`
	MaxMembers int                `bson:"maxMembers" json:"maxMembers"`
	ManagerID  primitive.ObjectID `bson:"managerId" json:"managerId"`
	IsActive   bool               `bson:"isActive" json:"isActive"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
}
