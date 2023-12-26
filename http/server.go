package http

import (
	"Blockride-waitlistAPI/internal/store"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	repository store.Service
	logger     *slog.Logger
}

func NewApplication(db *mongo.Collection, logger *slog.Logger) *Application {
	return &Application{
		repository: store.NewStore(db),
		logger:     logger,
	}
}
