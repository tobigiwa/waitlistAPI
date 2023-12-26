package store

type User struct {
	Name          string `json:"name" validate:"required" bson:"name"`
	Country       string `json:"country" validate:"required" bson:"country"`
	Email         string `json:"email" validate:"required,email" bson:"email"`
	SplWalletAddr string `json:"waddr" validate:"required,isOncurve" bson:"waddr"`
}
