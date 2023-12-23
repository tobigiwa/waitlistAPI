package store

type User struct {
	Name          string `json:"name" validate:"required" bson:"name"`
	Country       string `json:"country" validate:"required" bson:"country"`
	Email         string `json:"email" validate:"required,email" bson:"email"`
	SplWalletAddr string `json:"waddr" validate:"isOncurve" bson:"waddr"`
}

type CachedUser struct {
	RedisKeyStr string `json:"key"`
	U           User
}
