package http

import (
	"fmt"
	"net/http"
)

var (
	ErrLinkExpired  = fmt.Errorf("link expired")
	ErrDuplicateKey = fmt.Errorf("user with email account already exit")
)
func (a Application) clientError(w http.ResponseWriter, errStatus int, err error) {
	http.Error(w, err.Error(), errStatus)
}

func (a Application) serverError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
