package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	db    *mongo.Client
	cache *redis.Client
}

func NewStore(db *mongo.Client, rdb *redis.Client) Service {
	return Service{
		db:    db,
		cache: rdb,
	}
}

func (s Service) SetcacheWithExpiration(key string, value CachedUser) error {
	var (
		p   []byte
		err error
	)

	if p, err = json.Marshal(value); err != nil {
		return err
	}

	if err := s.cache.Set(context.Background(), key, p, 10*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}

func (s Service) GetFromCache(key string) (User, error) {

	var (
		cU     CachedUser
		result User
	)

	val, err := s.cache.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return result, NonExistentKey
	} else if err != nil {
		return result, err
	}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return result, err
	}

	return cU.U, nil
}

func (s Service) SaveToDb(user User) error {

	if _, err := s.db.Database("waitlist").Collection("user").InsertOne(context.TODO(), user); err != nil {
		if writeException, ok := err.(mongo.WriteException); ok {
			for _, writeError := range writeException.WriteErrors {
				if writeError.Code == 11000 {
					return errors.New("duplicate key error")
				}
			}
		}
		return err
	}

	return nil
}

func (s Service) CheckIfUserExist(user User) bool {
	var result User

	if err := s.db.Database("waitlist").Collection("user").FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&result); err == nil {
		return user.Email == result.Email // user exist

	} else if err == mongo.ErrNoDocuments {
		return false // user does not exist
	}
	// else just means an error occured, buh i detest "else" statement so...omitted

	return false // an error probably occurred, we assume user does not exist
}

func (s Service) DeleteFromCache(key string) error {
	if err := s.cache.Del(context.TODO(), key).Err(); err != nil {
		return err
	}
	return nil
}
