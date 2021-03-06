// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package core

import (
	"edgex-club/internal/authorization"
	"edgex-club/internal/config"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const (
	referer_HTTPS string = "https://www.edgexfoundry.club"
	referer_HTTP  string = "http://www.edgexfoundry.club"
)

func GeneralFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/login.html" || path == "/redirect.html" {
			http.FileServer(http.Dir("static")).ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(path, "/public/") {
			http.FileServer(http.Dir("static")).ServeHTTP(w, r)
			return
		}

		//Authorization
		var token string
		if tokenCookie, err := r.Cookie("Authorization"); err == nil {
			token = tokenCookie.Value
		}

		var claims *authorization.CustomClaims
		var ok bool

		claims, ok = authorization.CheckToken(token)
		//检测认证API是否携带有效jwt token
		if strings.HasPrefix(path, "/api/v1/auth") {
			if !ok {

				if config.Conf().Env().Prod {
					referer := r.Header.Get("Referer")
					if !strings.HasPrefix(referer, referer_HTTPS) || !strings.HasPrefix(referer, referer_HTTP) {
						http.Error(w, "Unauthorized.", http.StatusUnauthorized)
						return
					}
				}

				//ajax 请求
				if r.Header.Get("x-requested-with") != "" {
					http.Error(w, "Unauthorized.", http.StatusUnauthorized)
					return
				}
				http.Redirect(w, r, "/login.html", http.StatusPermanentRedirect)
				return
			}
		}
		if ok {
			credsUser, err := json.Marshal(claims.Credentials)
			if err != nil {
				log.Println("转换creds失败！")
				//ajax 请求
				if r.Header.Get("x-requested-with") != "" {
					http.Error(w, "Unauthorized.", http.StatusUnauthorized)
					return
				}
				http.Redirect(w, r, "/login.html", http.StatusPermanentRedirect)
				return
			}
			//重写请求头，用于下游服务中使用，比如一些handler中需要用到用户信息
			r.Header.Set("CredsUser", string(credsUser))
		}

		next.ServeHTTP(w, r)
	})
}
