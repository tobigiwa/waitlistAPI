package http

import (
	"Blockride-waitlistAPI/internal/store"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	repository store.Service
	logger     *slog.Logger
}

func NewApplication(db *mongo.Collection, rdb *redis.Client, logger *slog.Logger) *Application {
	return &Application{
		repository: store.NewStore(db, rdb),
		logger:     logger,
	}
}


