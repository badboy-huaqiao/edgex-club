// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package handler

import (
	"edgex-club/internal/core"
	"edgex-club/internal/model"
	repo "edgex-club/internal/repository"
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
	articles := repo.ArticleRepos.FindAllArticles()
	hotAuthors := repo.ArticleRepos.HotAuthor()
	hotArticles := repo.ArticleRepos.HotArticle()
	data := struct {
		Articles    []model.Article
		HotAuthors  []model.User
		HotArticles []model.Article
	}{
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
	article := repo.ArticleRepos.FindOne(userName, articleId)
	user := repo.UserRepos.FindOneByName(userName)
	articleCount := repo.ArticleRepos.UserArticleCount(userName)
	comments := repo.ArticleRepos.FindAllCommentByArticleId(articleId)
	hotArticles := repo.ArticleRepos.HotArticle()
	replysMap := make(map[string][]model.Reply)
	for _, c := range comments {
		replys := repo.ArticleRepos.FindAllReplyByCommentId(c.Id.Hex())
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
	a = repo.ArticleRepos.FindOne(userName, articleId)
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
	articleCount := repo.ArticleRepos.UserArticleCount(userName)

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
