package store

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Service struct {
	db *mongo.Collection
}

func NewStore(db *mongo.Collection) Service {
	return Service{
		db: db,
	}
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
