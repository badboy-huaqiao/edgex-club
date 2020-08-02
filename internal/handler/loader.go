// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package handler

import (
	"edgex-club/internal/core"
	"edgex-club/internal/model"
	repo "edgex-club/internal/repository"
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
	creds := genCredsUser(r)

	articles, _ := repo.ArticleRepositotyClient().FetchAll()
	hotAuthors, _ := repo.ArticleRepositotyClient().HotAuthor()
	hotArticles, _ := repo.ArticleRepositotyClient().HotArticle()
	data := struct {
		CredsUser   *model.Credentials
		Articles    []model.Article
		HotAuthors  []model.User
		HotArticles []model.Article
	}{
		CredsUser:   creds,
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

	creds := genCredsUser(r)

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
		CredsUser       *model.Credentials
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
		CredsUser:       creds,
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

	creds := genCredsUser(r)

	data := struct {
		ArticleId    string
		MD           string
		ArticleTitle string
		Type         string
		CredsUser    *model.Credentials
	}{
		ArticleId:    articleId,
		MD:           a.Content,
		ArticleTitle: a.Title,
		Type:         a.Type,
		CredsUser:    creds,
	}
	renderTemplate(w, "article_edit", "base", data)
}

func LoadArticleAddPage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	creds := genCredsUser(r)
	data := struct {
		CredsUser *model.Credentials
	}{
		CredsUser: creds,
	}

	renderTemplate(w, "article_add", "base", data)
}

func LoadUserHomePage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	userName := vars["userName"]
	creds := genCredsUser(r)

	var articles []model.Article
	var err error
	var filter, self bool
	var msgs []model.Message
	if creds.Name != userName {
		filter = true
	} else {
		if msgs, err = repo.MessageRepositotyClient().FetchAllByUserName(creds.Name); err != nil {
			http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
			return
		}
		self = true
	}

	u := repo.UserRepos.FindOneByName(userName)
	//TODO: 数量应该按照是否认证后显示已审核的数量
	articleCount, _ := repo.ArticleRepositotyClient().UserArticleCount(userName)

	if articles, err = repo.ArticleRepositotyClient().FindAllArticlesByUser(u.Id.Hex(), filter); err != nil {
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	data := struct {
		CredsUser    *model.Credentials
		Self         bool
		UserId       string
		UserName     string
		AvatarUrl    string
		ArticleCount int
		Articles     []model.Article
		Messages     []model.Message
	}{
		CredsUser:    creds,
		Self:         self,
		UserId:       u.Id.Hex(),
		UserName:     userName,
		AvatarUrl:    u.AvatarUrl,
		ArticleCount: articleCount,
		Articles:     articles,
		Messages:     msgs,
	}

	renderTemplate(w, "userhome", "base", data)
}
