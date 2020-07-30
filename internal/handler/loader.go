// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package handler

import (
	"edgex-club/internal/core"
	"edgex-club/internal/model"
	repo "edgex-club/internal/repository"
	"encoding/json"
	"fmt"
	"net/http"

	mux "github.com/gorilla/mux"
)

func renderTemplate(w http.ResponseWriter, name string, template string, data interface{}) {
	t, ok := core.TemplateStore[name]
	if !ok {
		http.Error(w, "template resource not found", http.StatusInternalServerError)
		return
	}
	err := t.ExecuteTemplate(w, template, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoadIndexPage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userInfo := r.Header.Get("inner-user")
	var credUser *model.Credentials
	if userInfo != "" {
		if err := json.Unmarshal([]byte(userInfo), credUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	articles, _ := repo.ArticleRepositotyClient().FetchAll()
	hotAuthors, _ := repo.ArticleRepositotyClient().HotAuthor()
	hotArticles, _ := repo.ArticleRepositotyClient().HotArticle()
	data := struct {
		CredUser    *model.Credentials
		Articles    []model.Article
		HotAuthors  []model.User
		HotArticles []model.Article
	}{
		CredUser:    credUser,
		Articles:    articles,
		HotAuthors:  hotAuthors,
		HotArticles: hotArticles,
	}
	renderTemplate(w, "index", "base", data)
}

func LoadArticlePage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)

	userName := vars["userName"]
	articleId := vars["articleId"]
	fmt.Printf("userName: %s\n", vars["userName"])
	article, _ := repo.ArticleRepositotyClient().FindOne(userName, articleId)
	user := repo.UserRepos.FindOneByName(userName)
	articleCount, _ := repo.ArticleRepositotyClient().UserArticleCount(userName)
	comments, _ := repo.CommentRepositotyClient().FindAllCommentByArticleId(articleId)
	hotArticles, _ := repo.ArticleRepositotyClient().HotArticle()
	replysMap := make(map[string][]model.Reply)
	for _, c := range comments {
		replys, _ := repo.ReplyRepositotyClient().FindAllReplyByCommentId(c.Id.Hex())
		replysMap[c.Id.Hex()] = replys
	}
	data := struct {
		UserName        string
		AvatarUrl       string
		ArticleId       string
		ArticleTitle    string
		ArticleCount    int
		ArticleModified int64
		ReadCount       int64
		MD              string
		Comments        []model.Comment
		ReplysMap       map[string][]model.Reply
		HotArticles     []model.Article
	}{
		UserName:        userName,
		AvatarUrl:       user.AvatarUrl,
		ArticleId:       articleId,
		ArticleTitle:    article.Title,
		ArticleCount:    articleCount,
		ArticleModified: article.Modified,
		ReadCount:       article.ReadCount,
		MD:              article.Content,
		Comments:        comments,
		ReplysMap:       replysMap,
		HotArticles:     hotArticles,
	}
	renderTemplate(w, "article", "base", data)
}

func LoadArticleEditPage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	articleId := vars["articleId"]
	userName := vars["userName"]
	var a model.Article
	a, _ = repo.ArticleRepositotyClient().FindOne(userName, articleId)
	data := struct {
		ArticleId    string
		MD           string
		ArticleTitle string
		Type         string
	}{
		ArticleId:    articleId,
		MD:           a.Content,
		ArticleTitle: a.Title,
		Type:         a.Type,
	}
	renderTemplate(w, "article_edit", "base", data)
}

func LoadArticleAddPage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	renderTemplate(w, "article_add", "base", nil)
}

func LoadUserHomePage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	userName := vars["userName"]
	u := repo.UserRepos.FindOneByName(userName)
	articleCount, _ := repo.ArticleRepositotyClient().UserArticleCount(userName)

	data := struct {
		UserId       string
		UserName     string
		AvatarUrl    string
		ArticleCount int
	}{
		UserId:       u.Id.Hex(),
		UserName:     userName,
		AvatarUrl:    u.AvatarUrl,
		ArticleCount: articleCount,
	}

	renderTemplate(w, "userhome", "base", data)
}
