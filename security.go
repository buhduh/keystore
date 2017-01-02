package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	b32 "encoding/base32"
	b64 "encoding/base64"
	"fmt"
	"github.com/balasanjay/totp"
	"golang.org/x/crypto/bcrypt"
	"io"
	"keystore/session"
	"net/http"
)

const (
	//This will likely need to be
	//tuned given the power of the pi
	BCRYPT_COST int = 12
)

type Verifier interface {
	//Does the request require security?
	IsSecure(*http.Request) bool
	IsLoggedIn(*http.Request) bool
}

type TwoWayDecryptor interface {
	EncryptPassword(string, string) (string, error)
	DecryptPassword(string, string) (string, error)
}

type TwoWayDecrypt struct {
	key string
}

func NewTwoWay(key string) *TwoWayDecrypt {
	return &TwoWayDecrypt{key: key}
}

type DefaultVerifier struct {
	secure bool
}

func (d DefaultVerifier) IsSecure(r *http.Request) bool {
	return d.secure
}

func (d DefaultVerifier) IsLoggedIn(r *http.Request) bool {
	sess, err := session.NewSession(r.Cookies(), session.DEFAULT_SESSION_AGE)
	if err != nil || sess == nil {
		return false
	}
	return sess.Get("logged_in") == "true"
}

type FormVerifier struct{}

//ALL submitted forms must be scrutinized
func (f FormVerifier) IsSecure(r *http.Request) bool {
	return true
}

//Slightly misnamed, but if a form does not require a user to be logged in,
//this should return true to satisfy Route.ServeHTTP()
func (f FormVerifier) IsLoggedIn(r *http.Request) bool {
	nonce := r.FormValue(NONCE_FORM_NAME)
	if nonce == "" {
		return false
	}
	validNonce := Nonce(nonce).CheckNonce()
	if !validNonce {
		return false
	}
	reqLogin := r.FormValue(REQ_LOGIN_FORM_NAME)
	if reqLogin == "" {
		return false
	}
	if reqLogin != "false" && reqLogin != "true" {
		return false
	}
	if reqLogin == "false" {
		return true
	}
	//if there's no action, the form can't be properly parsed, bail
	action := r.FormValue(ACTION_FORM_NAME)
	if action == "" {
		return false
	}
	sess, err := session.NewSession(r.Cookies(), session.DEFAULT_SESSION_AGE)
	if err != nil || sess == nil {
		return false
	}
	return sess.Get("logged_in") == "true"
}

func scramble(salt, pw string) []byte {
	return []byte(fmt.Sprintf("%s%s%s%s", salt, pw, pw, salt))
}

func GetPassord(salt, pw string) (string, error) {
	p, err := bcrypt.GenerateFromPassword(scramble(salt, pw), BCRYPT_COST)
	if err != nil {
		return "", err
	}
	return string(p), nil
}

func ComparePassword(salt, pw, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), scramble(salt, pw))
	return err == nil
}

func VerifyCode(secret, code string) bool {
	bSecret, err := b32.StdEncoding.DecodeString(secret)
	if err != nil {
		return false
	}
	return totp.Authenticate(bSecret, code, nil)
}

func generateKey(salt, key string) []byte {
	temp := sha256.Sum256([]byte(salt + key + key + salt))
	return temp[:]
}

//returns base 64 encoded encrypted pw
func (t TwoWayDecrypt) EncryptPassword(salt, pw string) (string, error) {
	block, err := aes.NewCipher(generateKey(salt, t.key))
	if err != nil {
		return "", err
	}
	b := b64.StdEncoding.EncodeToString([]byte(pw))
	cipherText := make([]byte, aes.BlockSize+len(b))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], []byte(b))
	data := b64.StdEncoding.EncodeToString(cipherText)
	return data, nil
}

//takes base 64 encoded encrypted pw
func (t TwoWayDecrypt) DecryptPassword(salt, tPW string) (string, error) {
	block, err := aes.NewCipher(generateKey(salt, t.key))
	if err != nil {
		return "", err
	}
	//text := []byte(pw)
	text, err := b64.StdEncoding.DecodeString(tPW)
	if err != nil {
		return "", err
	}
	if len(text) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := b64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
