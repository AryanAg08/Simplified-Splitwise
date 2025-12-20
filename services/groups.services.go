package services

import (
	"context"
	"errors"
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

	group.ID = result.InsertedID.(primitive.ObjectID)

	_, balanceError := InitBalance(group.ID.Hex(), members)
	if balanceError != nil {
		return models.Groups{}, balanceError
	}

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

	balanceError := AddMembersToBalance(groupId, members)
	if balanceError != nil {
		return models.Groups{}, balanceError
	}

	return updatedGroup, nil
}

func (g *GroupsService) GroupDetails(groupId string) (models.Groups, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, errr := primitive.ObjectIDFromHex(groupId)

	if errr != nil {
		return models.Groups{}, errors.New("invalid group id")
	}

	var group models.Groups

	err := db.GroupCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&group)

	if err != nil {
		return models.Groups{}, err
	}

	return group, nil
}

func (g *GroupsService) GetAllGroups() ([]models.Groups, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var groups []models.Groups

	cursor, err := db.GroupCollection.Find(ctx, bson.M{})

	if err != nil {
		return []models.Groups{}, err
	}

	if err := cursor.All(ctx, &groups); err != nil {
		return []models.Groups{}, err
	}

	return groups, nil
}

func InitBalance(groupID string, members []string) (bool, error) {
	if len(members) == 0 {
		return false, errors.New("no members found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	balances := make(map[string]float64)
	for _, m := range members {
		balances[m] = 0
	}

	balanceGroup := models.Balance{
		GroupID:   groupID,
		Balances:  balances,
		UpdatedAt: time.Now(),
	}

	_, err := db.BalanceCollection.InsertOne(ctx, balanceGroup)
	if err != nil {
		return false, err
	}

	return true, nil
}

func AddMembersToBalance(groupID string, members []string) error {
	if len(members) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	setFields := bson.M{}
	for _, m := range members {
		setFields["balances."+m] = 0
	}

	setFields["updatedAt"] = time.Now()

	update := bson.M{
		"$set": setFields,
	}

	res, err := db.BalanceCollection.UpdateOne(
		ctx,
		bson.M{"groupId": groupID},
		update,
	)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("balance document not found for group")
	}

	return nil
}
