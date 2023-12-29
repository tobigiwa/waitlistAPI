package http

import (
	"Blockride-waitlistAPI/env"
	"Blockride-waitlistAPI/internal/store"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

// WaitListHandler
//
//	@Summary		sends user registration email
//	@Description	sends user registration email
//	@Tags			application
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200				{string}	string	"OK"
//	@Param			name			formData	string	true	"Users any preferred name"
//	@Param			email			formData	string	true	"valid email address"	Format(email)
//	@Param			country			formData	string	true	"country"
//	@Param			splWalletAddr	formData	string	true	"SPL wallet Address"
//	@Failure		400				{string}	string	"CLIENT ERROR: BAD REQUEST, INVALID USER FORM DATA"
//	@Failure		409				{string}	string	"CLIENT ERROR: USER WITH EMAIL ALREADY EXIST"
//	@Failure		500				{string}	string	"SERVER ERROR: INTERNAL SERVRER ERROR"
//	@Router			/joinwaitlist [post]
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
	if name := r.PostForm.Get("name"); name != "" {
		subscriber.Name = name
	}
	if email := r.PostForm.Get("email"); email != "" {
		subscriber.Email = email
	}
	if country := r.PostForm.Get("country"); country != "" {
		subscriber.Country = strings.ToUpper(country)
	}
	if splWalletAddr := r.PostForm.Get("splWalletAddr"); splWalletAddr != "" {
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
		w.Header().Set("Location", "https://www.blockride.xyz/")
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

	w.Header().Set("Location", "https://www.blockride.xyz/")
	http.Redirect(w, r, "https://www.blockride.xyz/", http.StatusOK)
}

// ConfirmAndSaveHandler
//
//	@Summary		confirms user registration
//	@Description	confirms user registration from email link
//	@Tags			application
//	@Produce		json
//	@Param			k	query		string	true	"BASE64 ENCODED STRING"
//	@Success		200	{string}	string	"REDIRECT TO BLOCKRIDE HOMEPAGE"
//	@Failure		400	{string}	string	"CLIENT ERROR: BAD REQUEST, KEY MISSING IN REQUEST"
//	@Failure		404	{string}	string	"CLIENT ERROR: NOT FOUND, LINK/KEY EXPIRED"
//	@Failure		409	{string}	string	"CLIENT ERROR: USER WITH EMAIL ALREADY EXIST"
//	@Failure		500	{string}	string	"SERVER ERROR: INTERNAL SERVRER ERROR"
//	@Router			/confirmuser [get]
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
			a.clientError(w, http.StatusNotFound, ErrLinkExpired)
			return
		}
		a.serverError(w)
		a.logger.LogAttrs(context.TODO(), slog.LevelError, err.Error())
		return
	}

	if err := a.repository.SaveToDb(user); err != nil {
		if strings.Contains(err.Error(), "duplicate key error") {
			w.Header().Set("Location", "https://www.blockride.xyz/")
			http.Redirect(w, r, "https://www.blockride.xyz/", http.StatusConflict)
			return
		}
		
		a.serverError(w)
		a.logger.LogAttrs(context.TODO(), slog.LevelError, err.Error())
		return
	}

	w.Header().Set("Location", "https://www.blockride.xyz/")
	http.Redirect(w, r, "https://www.blockride.xyz/", http.StatusOK)
}

type ServerStatus struct {
	Server_status       string
	Application_Env     string
	Application_Version string
}

// HealthcheckHandler
//
//	@Summary		Report application status
//	@Description	return application status
//	@Tags			status
//	@Produce		json
//	@Success		200	{object}	http.ServerStatus	"Server_status:available"
//	@Failure		500	{string}	string				"INTERNAL SERVRER ERROR"
//	@Router			/healthcheck [get]
func (a Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	env := ServerStatus{
		Server_status:       "available",
		Application_Env:     env.GetEnvVar().Server.Env,
		Application_Version: env.GetEnvVar().Server.Version,
	}
	var (
		byteArr []byte
		err     error
	)
	if byteArr, err = json.Marshal(env); err != nil {
		a.logger.LogAttrs(context.TODO(), slog.LevelError, "marshling error from healthcheckHandler"+err.Error())
		a.serverError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(byteArr)
}
