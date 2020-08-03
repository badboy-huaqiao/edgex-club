// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package handler

import (
	"edgex-club/internal/model"
	repo "edgex-club/internal/repository"
	"html/template"
	"net/http"

	mux "github.com/gorilla/mux"
)

type TodoPageUserData struct {
	UserId       string
	UserName     string
	AvatarUrl    string
	ArticleCount int
	Articles     []model.Article
}

func UserHome(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	userName := vars["userName"]
	creds := genCredsUser(r)

	var articles []model.Article
	var err error
	var filter bool
	if creds.Name != userName {
		filter = true
	}

	if articles, err = repo.ArticleRepositotyClient().FindAllArticlesByUser(creds.Id, filter); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u, _ := repo.UserRepositoryClient().FetchOneByName(userName)

	articleCount, _ := repo.ArticleRepositotyClient().UserArticleCount(userName)

	t, _ := template.ParseFiles("static/users/home.html")
	data := TodoPageUserData{
		UserId:       u.Id.Hex(),
		UserName:     userName,
		AvatarUrl:    u.AvatarUrl,
		ArticleCount: articleCount,
		Articles:     articles,
	}
	t.Execute(w, data)
}
