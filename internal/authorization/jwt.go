// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package authorization

import (
	"edgex-club/internal/config"
	"edgex-club/internal/model"
	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("AllYourBase")

type CustomClaims struct {
	model.Credentials `json:"credentials"`
	jwt.StandardClaims
}

var (
	verifyKey, signKey []byte
)

func InitAuth() {
	var err error

	signKey, err = ioutil.ReadFile(config.Config.Service.JWTPrivKey)
	if err != nil {
		log.Fatal("Error reading jwt rsa private key")
		return
	}
	verifyKey, err = ioutil.ReadFile(config.Config.Service.JWTPubKey)
	if err != nil {
		log.Fatal("Error reading jwt rsa public key")
		return
	}
}

func NewToken(creds model.Credentials) (token string, err error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claim := &CustomClaims{
		Credentials: creds,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	token, err = jwtToken.SignedString(mySigningKey)
	if err != nil {
		log.Printf("创建JWT token签名失败：%v", err.Error())
		return "", err
	}
	return token, err
}

func CheckToken(token string) (claims *CustomClaims, ok bool) {
	claims = &CustomClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	//包括超时、被篡改等，都会无效
	if err != nil || !jwtToken.Valid {
		return nil, false
	}

	if claims, ok := jwtToken.Claims.(*CustomClaims); ok && jwtToken.Valid {
		return claims, true
	}

	return nil, false
}

func RefreshToken(tokenStr string) (string, error) {
	claims := &CustomClaims{}
	jwtToken, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil || !jwtToken.Valid {
		return "", err
	}

	//确保新的token发放的时间，不会在旧的token超期了才发放。
	//因此这里确保，在旧token还有30秒以上的时间才会过期时，才会创建新的token。否则返回错误。
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return "", nil
	}

	return NewToken(claims.Credentials)
}
