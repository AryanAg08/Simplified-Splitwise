package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Expense struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	GroupID      string             `bson:"groupId" json:"groupId"`
	Description  string             `bson:"description" json:"description"`
	PaidBy       string             `bson:"paidBy" json:"paidBy"`
	Amount       float64            `bson:"amount" json:"amount"`
	SplitBetween []string           `bson:"splitBetween" json:"splitBetween"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
}
