package models

type Groups struct {
	ID      string   `bson:"_id" json:"id"`
	Name    string   `bson:"name" json:"name"`
	Members []string `bson:"members" json:"members"`
}
