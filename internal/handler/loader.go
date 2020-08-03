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

type userInfo struct {
	Id        string
	Name      string
	AvatarUrl string
}

type pageDataPayload struct {
	IsSelf          bool
	CredsUser       *model.Credentials
	UserMsgSum      int
	UserArticlesSum int
	Article         model.Article
	Articles        []model.Article
	HotAuthors      []model.User
	HotArticles     []model.Article
	Comments        []model.Comment
	ReplysMap       map[string][]model.Reply
	Messages        []model.Message
	UserInof        userInfo
}

func renderTemplate(w http.ResponseWriter, name string, template string, data *pageDataPayload, creds *model.Credentials) {
	t, ok := core.TemplateStore[name]
	if !ok {
		http.Error(w, "template resource not found", http.StatusInternalServerError)
		return
	}
	if creds != nil {
		msgSum, _ := repo.MessageRepositotyClient().MsgSumByUserName(creds.Name)
		data.UserMsgSum = msgSum
		data.CredsUser = creds
	}
	err := t.ExecuteTemplate(w, template, data)
	if err != nil {
		fmt.Printf("exec render template error: %s\n", err.Error())
	}
}

func LoadIndexPage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	articles, _ := repo.ArticleRepositotyClient().FetchPageAll(0, 2)
	hotAuthors, _ := repo.ArticleRepositotyClient().HotAuthor()
	hotArticles, _ := repo.ArticleRepositotyClient().HotArticle()
	data := &pageDataPayload{
		Articles:    articles,
		HotAuthors:  hotAuthors,
		HotArticles: hotArticles,
	}
	renderTemplate(w, "index", "base", data, genCredsUser(r))
}

func LoadArticlePage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	userName := vars["userName"]
	articleId := vars["articleId"]

	article, _ := repo.ArticleRepositotyClient().FindOne(userName, articleId)
	articleCount, _ := repo.ArticleRepositotyClient().UserArticleCount(userName)
	comments, _ := repo.CommentRepositotyClient().FindAllCommentByArticleId(articleId)

	hotArticles, _ := repo.ArticleRepositotyClient().HotArticle()
	replysMap := make(map[string][]model.Reply)
	for _, c := range comments {
		replys, _ := repo.ReplyRepositotyClient().FindAllReplyByCommentId(c.Id.Hex())
		replysMap[c.Id.Hex()] = replys
	}
	data := &pageDataPayload{
		Article:         article,
		HotArticles:     hotArticles,
		Comments:        comments,
		ReplysMap:       replysMap,
		UserArticlesSum: articleCount,
	}
	renderTemplate(w, "article", "base", data, genCredsUser(r))
}

func LoadArticleEditPage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	articleId := vars["articleId"]
	credsUser := genCredsUser(r)
	a, _ := repo.ArticleRepositotyClient().FindOne(credsUser.Name, articleId)
	data := &pageDataPayload{Article: a}
	renderTemplate(w, "article_edit", "base", data, credsUser)
}

func LoadArticleAddPage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := &pageDataPayload{}
	renderTemplate(w, "article_add", "base", data, genCredsUser(r))
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

	if creds != nil && creds.Name == userName {
		self = true
		if msgs, err = repo.MessageRepositotyClient().FetchAllByUserName(creds.Name); err != nil {
			http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
			return
		}
	} else {
		filter = true
	}

	u, _ := repo.UserRepositoryClient().FetchOneByName(userName)
	//TODO: 数量应该按照是否认证后显示已审核的数量
	articleCount, _ := repo.ArticleRepositotyClient().UserArticleCount(userName)

	if articles, err = repo.ArticleRepositotyClient().FindAllArticlesByUser(u.Id.Hex(), filter); err != nil {
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	data := &pageDataPayload{
		UserArticlesSum: articleCount,
		Articles:        articles,
		Messages:        msgs,
		IsSelf:          self,
		UserInof: userInfo{
			Id:        u.Id.Hex(),
			Name:      userName,
			AvatarUrl: u.AvatarUrl,
		},
	}

	renderTemplate(w, "userhome", "base", data, creds)
}
