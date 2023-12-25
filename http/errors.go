package http

import (
	"net/http"
)

const (
	invalidUserDataError   = "User validation failed"
	unmarshalError         = "Json unmarshaling Error"
	marshalError           = "Json marshaling Error"
	formParsingError       = "Error parsing request form data"
	rediSetError           = "Error setting value to redis"
	rediGetError           = "Error geting value to redis"
	mailNotSentError       = "Error sending mail"
	unrecognisedKey        = "unrecognisedKey/invalid url"
	setToMongoDbError      = "Error saving to mongoDB"
	deleteFromMongoDbError = "Error deleting from mongoDB"
)

func (a Application) clientError(w http.ResponseWriter, errStatus int, err error) {
	http.Error(w, err.Error(), errStatus)
}

func (a Application) serverError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
