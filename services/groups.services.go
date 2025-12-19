package services

import (
	"context"
	"errors"

	// "log"
	"time"

	"github.com/AryanAg08/Simplified-Splitwise/models"
	"github.com/AryanAg08/Simplified-Splitwise/workers/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GroupsService struct{}

// Create a new group!!
// --> Remaining to error check and verify body!!
func (g *GroupsService) CreateGroupService(
	name string,
	members []string,
) (models.Groups, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	group := models.Groups{
		Name:    name,
		Members: members,
	}

	result, err := db.GroupCollection.InsertOne(ctx, group)
	if err != nil {
		return models.Groups{}, err
	}

	// 🔥 Assign MongoDB generated _id back to struct
	group.ID = result.InsertedID.(primitive.ObjectID)

	return group, nil
}

func (g *GroupsService) AddGroupMembers(
	groupId string,
	members []string,
) (models.Groups, error) {

	if len(members) == 0 {
		return models.Groups{}, errors.New("no members provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		return models.Groups{}, errors.New("invalid group id")
	}

	filter := bson.M{"_id": objID}

	update := bson.M{
		"$addToSet": bson.M{
			"members": bson.M{
				"$each": members,
			},
		},
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var updatedGroup models.Groups

	err = db.GroupCollection.
		FindOneAndUpdate(ctx, filter, update, opts).
		Decode(&updatedGroup)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Groups{}, errors.New("group not found")
		}
		return models.Groups{}, err
	}

	return updatedGroup, nil
}
