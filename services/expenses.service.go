package services

import (
	"context"
	"errors"
	"time"

	"github.com/AryanAg08/Simplified-Splitwise/models"
	"github.com/AryanAg08/Simplified-Splitwise/workers/db"
	"github.com/AryanAg08/Simplified-Splitwise/workers/queue"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExpensesSerive struct{}

func (e *ExpensesSerive) AddExpensesService(description string,
	paidBy string,
	amount float64,
	splitBetween []string,
	groupId string) (models.Expense, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	expense := models.Expense{
		GroupID:      groupId,
		Description:  description,
		PaidBy:       paidBy,
		Amount:       amount,
		SplitBetween: splitBetween,
		CreatedAt:    time.Now(),
	}
	// check if the grp exists and member exists ‼️‼️

	grpId, err := primitive.ObjectIDFromHex(groupId)

	if err != nil {
		return models.Expense{}, err
	}

	var group models.Groups

	err = db.GroupCollection.FindOne(ctx, bson.M{"_id": grpId}).Decode(&group)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Expense{}, errors.New("no such group exists")
		}
		return models.Expense{}, err
	}

	// check for user exists in that group!
	if !contains(group.Members, paidBy) {
		for _, v := range splitBetween {
			if !contains(group.Members, v) {
				return models.Expense{}, errors.New("split between user not found in group: " + v)
			}
		}
		return models.Expense{}, errors.New("paid by user not found in group: " + paidBy)
	}

	result, err := db.ExpenseCollection.InsertOne(ctx, expense)

	if err != nil {
		return models.Expense{}, err
	}

	expense.ID = result.InsertedID.(primitive.ObjectID)

	// send Message to Rabbit MQ to recalculate balance!!

	body := []byte(`{
		"groupId": "` + groupId + `",
		"expenseId": "` + expense.ID.Hex() + `"
	}`)

	queue.Channel.Publish(
		"",
		"expense_added",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	return expense, nil
}

func (e *ExpensesSerive) GetAllExpensesService(groupId string) ([]models.Expense, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var expenses []models.Expense

	cursor, err := db.ExpenseCollection.Find(ctx, bson.M{"GroupID": groupId})

	if err != nil {
		return []models.Expense{}, err
	}

	if err := cursor.All(ctx, &expenses); err != nil {
		return []models.Expense{}, err
	}

	return expenses, nil

}

func contains(arr []string, target string) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}
