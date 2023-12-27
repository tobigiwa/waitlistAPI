package http

import (
	"Blockride-waitlistAPI/env"
	"Blockride-waitlistAPI/internal/store"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func (a Application) waitListHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		a.clientError(w, http.StatusBadRequest, err)
		a.logger.LogAttrs(context.TODO(), slog.LevelError, err.Error())
		return
	}

	var (
		subscriber store.User
		keyStr     string
		err        error
	)
	if name := r.FormValue("name"); name != "" {
		subscriber.Name = name
	}
	if email := r.FormValue("email"); email != "" {
		subscriber.Email = email
	}
	if country := r.FormValue("country"); country != "" {
		subscriber.Country = strings.ToUpper(country)
	}
	if splWalletAddr := r.FormValue("splWalletAddr"); splWalletAddr != "" {
		subscriber.SplWalletAddr = splWalletAddr
	}

	if (subscriber == store.User{}) {
		a.clientError(w, http.StatusBadRequest, fmt.Errorf("empty form data"))
		a.logger.LogAttrs(context.TODO(), slog.LevelError, fmt.Errorf("empty form data").Error())
		return
	}

	if err := validateSubscriber(subscriber); err != nil {
		a.clientError(w, http.StatusBadRequest, fmt.Errorf("validation error: %w", err))
		a.logger.LogAttrs(context.TODO(), slog.LevelError, err.Error())
		return
	}

	if a.repository.CheckIfUserExist(subscriber) {
		w.Header().Set("Location", "https://www.blockride.xyz/") // this still need a different url/page
		http.Redirect(w, r, "https://www.blockride.xyz/", http.StatusConflict)
		return
	}

	if keyStr, err = encryptUserInfo(subscriber); err != nil {
		a.serverError(w)
		a.logger.LogAttrs(context.TODO(), slog.LevelError, err.Error())
		return
	}

	if err := sendConfirmationMail(subscriber.Name, subscriber.Email, keyStr); err != nil {
		a.serverError(w)
		a.logger.LogAttrs(context.TODO(), slog.LevelError, err.Error())
		return
	}

	w.Header().Set("Location", "https://www.blockride.xyz/") // this still need a different url/page
	http.Redirect(w, r, "https://www.blockride.xyz/", http.StatusOK)
}

func (a Application) confirmAndSaveHandler(w http.ResponseWriter, r *http.Request) {

	var (
		user   store.User
		err    error
		keyStr string
	)

	if keyStr = r.URL.Query().Get("k"); keyStr == "" {
		a.clientError(w, http.StatusBadRequest, errors.New("not a valid url"))
		a.logger.LogAttrs(context.TODO(), slog.LevelInfo, err.Error())
		return
	}

	if user, err = dencryptUserInfo(keyStr); err != nil {
		if errors.Is(err, ErrLinkExpired) {
			a.clientError(w, http.StatusBadRequest, ErrLinkExpired)
			return
		}
		a.serverError(w)
		a.logger.LogAttrs(context.TODO(), slog.LevelError, err.Error())
		return
	}

	if err := a.repository.SaveToDb(user); err != nil {
		if strings.Contains(err.Error(), "duplicate key error") {
			a.clientError(w, http.StatusBadRequest, ErrDuplicateKey)
			return
		}
		a.serverError(w)
		a.logger.LogAttrs(context.TODO(), slog.LevelError, err.Error())
		return
	}

	w.Header().Set("Location", "https://www.blockride.xyz/")
	http.Redirect(w, r, "https://www.blockride.xyz/", http.StatusSeeOther)
}

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	env := struct {
		server_status    string
		application_info map[string]string
	}{
		server_status: "available",
		application_info: map[string]string{
			"enviroment": env.GetEnvVar().Server.Env,
			"version":    env.GetEnvVar().Server.Version,
		},
	}

	w.Write([]byte(fmt.Sprintf("%+v", env)))
}
