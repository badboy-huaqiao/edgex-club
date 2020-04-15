// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package core

import (
	"edgex-club/internal/authorization"
	"edgex-club/internal/errors"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func GeneralFilter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasPrefix(path, "/public/") {
			http.FileServer(http.Dir("static")).ServeHTTP(w, r)
			return
		}
		//检测认证API是否携带有效jwt token
		if strings.HasPrefix(path, "/api/v1/auth") {
			token := r.Header.Get("edgex-club-token")
			if token == "" {
				tokenCookie, err := r.Cookie("edgex-club-token")
				if err != nil {
					http.Error(w, errors.NewErrUnauthorize().Error(), http.StatusUnauthorized)
					return
				}
				token = tokenCookie.Value
			}
			claims := &authorization.Claims{}
			ok := false
			if ok, claims = authorization.CheckToken(token); !ok {
				http.Error(w, errors.NewErrUnauthorize().Error(), http.StatusUnauthorized)
				return
			}
			credsByte, err := json.Marshal(claims.Credentials)

			if err != nil {
				log.Println("转换creds失败！")
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			//重写请求头，用于下游服务中使用，比如一些controller中需要用到用户信息
			r.Header.Set("inner-user", string(credsByte))
		}

		h.ServeHTTP(w, r)
	})
}
