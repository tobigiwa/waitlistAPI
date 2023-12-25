package http

import (
	"Blockride-waitlistAPI/internal/store"

	solana "github.com/blocto/solana-go-sdk/common"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func validateSubscriber(s store.User) error {

	validate = validator.New(validator.WithRequiredStructEnabled())

	if err := validate.RegisterValidation("isOncurve", isOncurve); err != nil {
		return err
	}

	if err := validate.Struct(s); err != nil {
		return err
	}

	return nil
}

func isOncurve(fl validator.FieldLevel) bool {
	pubKey := solana.PublicKeyFromString(fl.Field().String())
	return solana.IsOnCurve(pubKey)
}
