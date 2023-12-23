package http

import (
	"Blockride-waitlistAPI/internal/store"
	"errors"
	"fmt"
	"net/http"
)

func (a Application) waitListHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		clientError(w, http.StatusBadRequest, err)
		return
	}

	var subscriber store.User

	if name := r.PostForm.Get("name"); name != "" {
		subscriber.Name = name
	}
	if email := r.PostForm.Get("email"); email != "" {
		subscriber.Email = email
	}
	if country := r.PostForm.Get("country"); country != "" {
		subscriber.Country = country
	}
	if splWalletAddr := r.PostForm.Get("splWalletAddr"); splWalletAddr != "" {
		subscriber.SplWalletAddr = splWalletAddr
	}

	if err := validateSubscriber(subscriber); err != nil {
		clientError(w, invalidFormData, err)
		return
	}

	if a.repository.CheckIfUserExist(subscriber) {
		clientError(w, http.StatusConflict, errors.New("email already exists"))
		return
	}

	keyStr := generateRedisKey(fmt.Sprintf("%s:%s", subscriber.Name, subscriber.Email))

	cachedUser := store.CachedUser{
		RedisKeyStr: keyStr,
		U:           subscriber,
	}

	if err := a.repository.SetcacheWithExpiration(keyStr, cachedUser); err != nil {
		serverError(w, err)
		return
	}

	if err := sendConfirmationMail(subscriber.Name, subscriber.Email, keyStr); err != nil {
		serverError(w, err)
		return
	}

	fmt.Println("waitListHanler good")
}

func (a Application) confirmAndSaveHandler(w http.ResponseWriter, r *http.Request) {

	var (
		user   store.User
		err    error
		keyStr string
	)

	if keyStr := r.URL.Query().Get("k"); keyStr == "" {
		clientError(w, http.StatusBadRequest, errors.New("not a valid url"))
		return
	}

	if user, err = a.repository.GetFromCache(keyStr); err != nil {
		serverError(w, err)
		return
	}

	if err := a.repository.SaveToDb(user); err != nil {
		serverError(w, err)
		return
	}
	if err := a.repository.DeleteFromCache(keyStr); err != nil {
		serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusFound)

	w.Header().Set("Location", "https://www.blockride.xyz/")

	http.Redirect(w, r, "https://www.blockride.xyz/", http.StatusSeeOther)

	fmt.Println("confirmAndSaveHandler good")
}
