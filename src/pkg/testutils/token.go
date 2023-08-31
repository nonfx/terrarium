// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"time"

	"github.com/golang-jwt/jwt"
	jose "gopkg.in/square/go-jose.v2"
)

const (
	key        = "mysupersecretkey"
	TestUserID = "3106056c-bdc6-4f06-be33-5056ba745023"
)

func GetTestToken() string {
	token, _ := CreateToken(key, TestUserID, "testuser@cldcvr.com", "test", "test", time.Now())
	return token
}

func GetExpiredToken() string {
	token, _ := CreateToken(key, TestUserID, "test", "test", "test", time.Now().AddDate(0, 0, -1))
	return token
}

func GetTestRefreshToken() string {
	refreshToken, _ := CreateRefreshToken(TestUserID, key, time.Now())
	return refreshToken
}

type Claims struct {
	UserID        string `json:"userId"`
	Email         string `json:"email"`
	ProviderToken string `json:"providerToken"`
	GCtoken       string `json:"gcToken"`
	jwt.StandardClaims
}

const (
	expirationTime        = time.Hour * 12
	refreshExpirationTime = time.Hour * 720
)

func CreateToken(jwtKey, userID, email, gcToken, providerToken string, currentTime time.Time) (string, error) {
	key := []byte(jwtKey)
	pt, err := joseEncryption(providerToken, key)
	if err != nil {
		return "", err
	}
	claim := Claims{
		Email:         email,
		UserID:        userID,
		ProviderToken: pt,
		GCtoken:       gcToken,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentTime.Unix(),
			ExpiresAt: currentTime.Add(expirationTime).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	cipher, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return cipher, nil
}

func joseEncryption(jwtToken string, key []byte) (string, error) {
	rcpt := jose.Recipient{
		Algorithm: jose.A256GCMKW,
		Key:       key,
	}
	opts := &jose.EncrypterOptions{
		Compression: jose.DEFLATE,
	}
	encrypter, err := jose.NewEncrypter(jose.A128GCM, rcpt, opts)
	if err != nil {
		return "", err
	}
	enc, err := encrypter.Encrypt([]byte(jwtToken))
	if err != nil {
		return "", err
	}
	return enc.CompactSerialize()
}

func CreateRefreshToken(userID, jwtKey string, currentTime time.Time) (string, error) {
	key := []byte(jwtKey)
	claim := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentTime.Unix(),
			ExpiresAt: currentTime.Add(refreshExpirationTime).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	cipher, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return cipher, nil
}
