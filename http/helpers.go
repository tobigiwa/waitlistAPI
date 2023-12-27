package http

import (
	"Blockride-waitlistAPI/env"
	"Blockride-waitlistAPI/internal/store"
	"Blockride-waitlistAPI/templ"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"time"

	gomail "gopkg.in/mail.v2"
)

type templData struct {
	Name string
	Key  string
}

func sendConfirmationMail(name, email, key string) error {

	var (
		from     = env.GetEnvVar().Mail.EmailAcc
		tpl      bytes.Buffer
		s        = templData{Name: name, Key: key}
		data     []byte
		err      error
	)
	
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "BlockRide Confirmation mail")

	if err = prepareMail(s, &tpl); err != nil {
		return err
	}

	if data, err = templ.BlockRideLogo.ReadFile("BlockRideLogo.png"); err != nil {
		return err
	}
	m.Attach("BlockRideLogo.png", gomail.SetHeader(map[string][]string{
		"Content-ID": {"BlockRideLogo"},
	}), gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(data)
		return err
	}))

	m.SetBody("text/html", tpl.String())

	d := gomail.NewDialer(env.GetEnvVar().Mail.EmailSmtpServerHost, env.GetEnvVar().Mail.EmailSmtpServerPort, from, env.GetEnvVar().Mail.EmailPswd)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func prepareMail(data templData, mailBuf *bytes.Buffer) error {

	t, err := template.ParseFS(templ.EmailHTML, "mail.html")
	if err != nil {
		return err
	}

	if err := t.Execute(mailBuf, data); err != nil {
		return err
	}

	return nil
}

type encondedUser struct {
	Expiration time.Time
	U          store.User
}

func encryptUserInfo(data store.User) (string, error) {

	var (
		byteArr []byte
		err     error
	)

	u := encondedUser{
		Expiration: time.Now().Add(15 * time.Minute),
		U:          data,
	}

	if byteArr, err = json.Marshal(u); err != nil {
		return "", err
	}

	if byteArr, err = encrypt(byteArr); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(byteArr), nil

}

func dencryptUserInfo(s string) (store.User, error) {
	var (
		byteArr []byte
		err     error
		eU      encondedUser
		user    store.User
	)

	if byteArr, err = base64.URLEncoding.DecodeString(s); err != nil {
		return user, err
	}
	if byteArr, err = decrypt(byteArr); err != nil {
		return user, err
	}

	if err = json.Unmarshal(byteArr, &eU); err != nil {
		return user, err
	}

	if checkIfLinkNotExpired(eU.Expiration) {
		return user, ErrLinkExpired
	}

	return eU.U, nil
}

func encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(env.GetEnvVar().EncryptionKey))
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decrypt(ciphertext []byte) ([]byte, error) {

	block, err := aes.NewCipher([]byte(env.GetEnvVar().EncryptionKey))
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func checkIfLinkNotExpired(t time.Time) bool {
	return time.Now().After(t)
}
