package nodeWorker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/AryanAg08/Simplified-Splitwise/models"
	"github.com/AryanAg08/Simplified-Splitwise/workers/cache"
	"github.com/AryanAg08/Simplified-Splitwise/workers/db"
	"github.com/AryanAg08/Simplified-Splitwise/workers/queue"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExpenseEvent struct {
	GroupId   string `json:"groupId"`
	ExpenseId string `json:"expenseId"`
}

func ExpenseAddedWorker() {

	msgs, err := queue.Channel.Consume(
		"expense_added",
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for msg := range msgs {

			var event ExpenseEvent

			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Println("Invalid message body:", err)
				msg.Nack(false, false)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			expenseID, err := primitive.ObjectIDFromHex(event.ExpenseId)
			if err != nil {
				log.Println("Invalid expense ID:", err)
				cancel()
				msg.Nack(false, false)
				continue
			}

			var expense models.Expense
			err = db.ExpenseCollection.FindOne(
				ctx,
				bson.M{"_id": expenseID},
			).Decode(&expense)

			if err != nil {
				log.Println("Expense not found:", err)
				cancel()
				msg.Nack(false, true) // retry
				continue
			}

			splitCount := float64(len(expense.SplitBetween))
			if splitCount == 0 {
				cancel()
				msg.Ack(false)
				continue
			}

			splitAmount := expense.Amount / splitCount

			inc := bson.M{}

			for _, member := range expense.SplitBetween {
				inc["balances."+member] = -splitAmount
			}

			// inc["balances."+expense.PaidBy] = expense.Amount - splitAmount

			update := bson.M{
				"$inc": inc,
				"$set": bson.M{
					"updatedAt": time.Now(),
				},
			}

			// _, err = db.BalanceCollection.UpdateOne(
			// 	ctx,
			// 	bson.M{"groupId": event.GroupId},
			// 	update,
			// )
			opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

			var updatedBalance models.Balance

			err = db.BalanceCollection.FindOneAndUpdate(
				ctx,
				bson.M{"groupId": event.GroupId},
				update,
				opts,
			).Decode(&updatedBalance)

			// cache final balance
			CacheGroupBalance(event.GroupId, updatedBalance)

			cancel()

			if err != nil {
				log.Println("Balance update failed:", err)
				msg.Nack(false, true)
				cancel()
				continue
			}

			log.Println("✅ Balance updated for group:", event.GroupId)

			msg.Ack(false)
		}
	}()

}

func CacheGroupBalance(groupId string, balance models.Balance) {
	ctx, cancel := context.WithTimeout(cache.Ctx, 5*time.Second)
	defer cancel()

	key := "balance:" + groupId

	data, err := json.Marshal(balance)

	if err != nil {
		log.Fatal("redis marshal error")
		return
	}

	pipe := cache.RedisClient.Pipeline()

	pipe.Set(
		ctx,
		key,
		data,
		10*time.Minute,
	)

	pipe.Del(
		ctx,
		"balances:all",
	)

	_, err = pipe.Exec(ctx)

	if err != nil {
		log.Println("redis pipeline error")
		return
	}
}
