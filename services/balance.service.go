package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/AryanAg08/Simplified-Splitwise/models"
	"github.com/AryanAg08/Simplified-Splitwise/workers/cache"
	"github.com/AryanAg08/Simplified-Splitwise/workers/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BalanceService struct{}

func (b *BalanceService) GetGroupBalance(groupId string) (*models.Balance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cacheKey := "balance:" + groupId

	//checking redis cache!!
	cached, err := cache.RedisClient.Get(ctx, cacheKey).Result()

	if err == nil {
		var balance models.Balance
		if err := json.Unmarshal([]byte(cached), &balance); err == nil {
			return &balance, nil
		}
		// unmarshal fails?
	}

	var balance models.Balance

	err = db.BalanceCollection.FindOne(
		ctx,
		bson.M{"groupId": groupId},
	).Decode(&balance)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("balance with groupId not found")
		}
		return nil, err
	}

	go cacheBalance(cacheKey, balance)

	return &balance, nil

}

func cacheBalance(key string, balance models.Balance) {
	ctx, cancel := context.WithTimeout(cache.Ctx, 2*time.Second)
	defer cancel()

	data, err := json.Marshal(balance)
	if err != nil {
		log.Fatal("error cache marshal")
		return
	}

	err = cache.RedisClient.Set(
		ctx,
		key,
		data,
		10*time.Minute,
	).Err()

	if err != nil {
		log.Fatal("redis set failed")
		return
	}
}
