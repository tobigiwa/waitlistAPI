package main

import (
	"Blockride-waitlistAPI/env"
	blockride "Blockride-waitlistAPI/http"
	"Blockride-waitlistAPI/internal/store"
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"strings"

	"Blockride-waitlistAPI/docs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//	@title						BlockRide waitlist-backend
//	@description				BlockRide waitlist-backend API endpoints.
//	@x-logo						{"url": "https://example.com/img.png", "backgroundColor": "#000000", "altText": "example logo", "href": "https://example.com/img.png"}
//
// contact.name   BlockRide
//
//	@contact.url				https://www.blockride.xyz/
//
// contact.email  giwaoluwatobi@gmail.com
//
//	@externalDocs.description	OpenAPI
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
	env.SetDefaults()

	docs.SwaggerInfo.Version = env.GetEnvVar().Server.Version
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	log.Println("Server is configured to run in", strings.ToUpper(env.GetEnvVar().Server.Env), "mode, at version -", strings.ToUpper(env.GetEnvVar().Server.Version))

	// setup MongoDB connection
	var (
		mongoDbClient *mongo.Client
		err           error
		uri           string
	)

	uri = "mongodb://localhost:27017/"
	if env.GetEnvVar().Server.Env == blockride.Production {
		uri = env.BuildURI(
			env.GetEnvVar().Databases.Mongo.Username,
			env.GetEnvVar().Databases.Mongo.Password,
			env.GetEnvVar().Databases.Mongo.Host)
	}

	if mongoDbClient, err = store.NewMongoClient(uri); err != nil {
		log.Fatalf("could not setup MongoDB connection (uri - `%s`): %v", uri, err)
	}
	log.Println("MongoDB connection successful")

	mongoDbCollection := mongoDbClient.Database(env.GetEnvVar().Databases.Mongo.Database).Collection(env.GetEnvVar().Databases.Mongo.Collection)

	// index `email` key as unique
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := mongoDbCollection.Indexes().CreateOne(context.TODO(), indexModel); err != nil {
		log.Fatalf("could not index `email` field in MongoDB err: %v", err)
	}
	log.Println("MongoDB collection indexed by `email` field successfully")

	// setup Application Server
	app := blockride.NewApplication(mongoDbCollection, logger)

	srv := &http.Server{
		Addr:         ":" + env.GetEnvVar().Server.Port,
		Handler:      app.Routes(),
		ErrorLog:     slog.NewLogLogger(jsonLogger, slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done
		close(done)

		// Graceful shutdown
		if err := srv.Shutdown(context.TODO()); err != nil {
			log.Fatalf("Server Shutdown Failed:%+v", err)
		}
	}()

	log.Println("Server is Running on", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen and serve returned an err: %v", err)
	}

	// handles MongoDB disconnection
	if err = mongoDbClient.Disconnect(context.TODO()); err != nil {
		log.Fatalf("could not successfully discconnect MongoDB connections err: %v", err)
	}

	log.Println("MongoDB disconnection successful")

	log.Println("Server Stopped")

}
