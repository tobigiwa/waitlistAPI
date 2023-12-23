package store

import "errors"

var (
	NonExistentKey = errors.New("key does not exist")
)