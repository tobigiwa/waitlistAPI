package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Service struct {
	db    *mongo.Collection
	cache *redis.Client
}

func NewStore(db *mongo.Collection, rdb *redis.Client) Service {
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
		return result, errors.New("key does not exist")
	} else if err != nil {
		return result, err
	}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return result, err
	}

	return cU.U, nil
}

func (s Service) SaveToDb(user User) error {

	if _, err := s.db.InsertOne(context.TODO(), user); err != nil {
		return err
	}
	return nil
}

func (s Service) CheckIfUserExist(user User) bool {

	var result User
	if err := s.db.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&result); err == nil {
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

func NewMongoClient(url string) (*mongo.Client, error) {

	var (
		client *mongo.Client
		err    error
	)
	severAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(url).SetServerAPIOptions(severAPI)

	if client, err = mongo.Connect(context.TODO(), opts); err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}
