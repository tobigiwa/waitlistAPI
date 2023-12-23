package http

import (
	"Blockride-waitlistAPI/internal/store"
)

type Application struct {
	repository store.Service
}
