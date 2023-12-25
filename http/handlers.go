package http

import (
	"Blockride-waitlistAPI/env"
	"Blockride-waitlistAPI/internal/store"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

func (a Application) waitListHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		a.clientError(w, http.StatusBadRequest, err)
		a.Log(slog.LevelError, err, formParsingError)
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
		a.clientError(w, http.StatusBadRequest, err)
		a.Log(slog.LevelError, err, invalidUserDataError)
		return
	}

	if a.repository.CheckIfUserExist(subscriber) {
		w.Header().Set("Location", "https://www.blockride.xyz/") // this still need a different url/page
		http.Redirect(w, r, "https://www.blockride.xyz/", http.StatusConflict)
		return
	}

	keyStr := generateRedisKey(fmt.Sprintf("%s:%s:%s", subscriber.Name, env.GetEnvVar().Nonce, subscriber.Email))

	cachedUser := store.CachedUser{
		RedisKeyStr: keyStr,
		U:           subscriber,
	}

	if err := a.repository.SetcacheWithExpiration(keyStr, cachedUser); err != nil {
		a.serverError(w)
		a.Log(slog.LevelError, err, rediSetError)
		return
	}

	if err := sendConfirmationMail(subscriber.Name, subscriber.Email, keyStr); err != nil {
		a.serverError(w)
		a.Log(slog.LevelError, err, mailNotSentError)
		return
	}

	w.Header().Set("Location", "https://www.blockride.xyz/") // this still need a different url/page
	http.Redirect(w, r, "https://www.blockride.xyz/", http.StatusOK)

	log.Println("waitListHanler good")
}

func (a Application) confirmAndSaveHandler(w http.ResponseWriter, r *http.Request) {

	var (
		user   store.User
		err    error
		keyStr string
	)

	if keyStr := r.URL.Query().Get("k"); keyStr == "" {
		a.clientError(w, http.StatusBadRequest, errors.New("not a valid url"))
		a.Log(slog.LevelInfo, err, unrecognisedKey)
		return
	}

	if user, err = a.repository.GetFromCache(keyStr); err != nil {
		a.serverError(w)
		a.Log(slog.LevelError, err, rediGetError)
		return
	}

	if err := a.repository.SaveToDb(user); err != nil {
		a.serverError(w)
		a.Log(slog.LevelError, err, setToMongoDbError)
		return
	}
	if err := a.repository.DeleteFromCache(keyStr); err != nil {
		a.serverError(w)
		a.Log(slog.LevelError, err, deleteFromMongoDbError)
		return
	}

	w.Header().Set("Location", "https://www.blockride.xyz/")
	http.Redirect(w, r, "https://www.blockride.xyz/", http.StatusSeeOther)

	log.Println("confirmAndSaveHandler good")
}
