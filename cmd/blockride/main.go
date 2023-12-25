package main

import (
	"Blockride-waitlistAPI/env"
	blockride "Blockride-waitlistAPI/http"
	"Blockride-waitlistAPI/internal/store"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// logging services
	opts := slog.HandlerOptions{
		AddSource: true,
	}

	jsonLogger := slog.NewJSONHandler(os.Stdout, &opts)
	logger := slog.New(jsonLogger)
	slog.SetDefault(logger)

	// load.env variables
	env.LoadAllEnvVars()

	// setup MongoDB connection
	var (
		mongoDbClient *mongo.Client
		err           error
	)
	uri := env.BuildURI(
		env.GetEnvVar().Databases.Mongo.Username,
		env.GetEnvVar().Databases.Mongo.Password,
		env.GetEnvVar().Databases.Mongo.Host,
	)
	if mongoDbClient, err = store.NewMongoClient(uri); err != nil {
		log.Fatalf("could not setup MongoDB connection: %v", err)
	}
	log.Println("MongoDB connection successful")

	mongoDbCollection := mongoDbClient.Database(env.GetEnvVar().Databases.Mongo.Database).Collection(env.GetEnvVar().Databases.Mongo.Collection)

	// index `email` key as unique
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := mongoDbCollection.Indexes().CreateOne(context.TODO(), indexModel); err != nil {
		panic(fmt.Errorf("could not index `email` field in MongoDB err: %v", err))
	}

	defer func() {
		// handles MongoDB disconnection
		if err = mongoDbClient.Disconnect(context.TODO()); err != nil {
			panic(fmt.Errorf("could not successfully discconnect MongoDB connections err: %v", err))
		}
	}()

	// setup Redis connection
	redisOpts := &redis.Options{
		Addr:     env.GetEnvVar().Databases.Redis.Host,
		Username: env.GetEnvVar().Databases.Redis.User,
		Password: env.GetEnvVar().Databases.Redis.Password,
		DB:       env.GetEnvVar().Databases.Redis.Db,
	}
	rdb := redis.NewClient(redisOpts)
	log.Println("Redis connection successful")

	// setup Application Server
	app := blockride.NewApplication(mongoDbCollection, rdb, logger)

	srv := &http.Server{
		Addr:         ":" + env.GetEnvVar().PORT.HTTP,
		Handler:      app.Routes(),
		ErrorLog:     slog.NewLogLogger(jsonLogger, slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer close(done)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve returned err: %v", err)
		}
	}()

	<-done
	log.Println("Server Stopped")

	// Graceful shutdown
	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

}
