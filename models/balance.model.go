package models

import "time"

type Balance struct {
	GroupID   string             `bson:"groupId" json:"groupId"`
	Balances  map[string]float64 `bson:"balances" json:"balances"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
