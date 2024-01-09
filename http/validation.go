package http

import (
	"companyXYZ-waitlistAPI/internal/store"

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

// func isValidSolanaAddress(fl validator.FieldLevel) bool {
// 	var (
// 		bytes []byte
// 		err   error
// 	)
// 	if bytes, err = base58.Decode(fl.Field().String()); err != nil {
// 		return false
// 	}
// 	return len(bytes) == ed25519.PublicKeySize
// }

// func isValidSolanaAddress1(fl validator.FieldLevel) bool {

// 	var (
// 		// acc client.AccountInfo
// 		err error
// 	)
// 	c := client.NewClient(rpc.MainnetRPCEndpoint)
// 	if _, err = c.GetAccountInfo(context.TODO(), fl.Field().String()); err != nil {
// 		return false
// 	}
// 	return true
// }
