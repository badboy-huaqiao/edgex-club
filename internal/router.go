// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package internal

import (
	"edgex-club/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func InitRouter() http.Handler {
	r := mux.NewRouter()

	//加载首页
	r.HandleFunc("/", handler.LoadIndexPage)

	s := r.PathPrefix("/api/v1").Subrouter()

	//+++++++++++++++++认证API++++++++++++++++++++++++++++++++++++++++++++++++

	s.HandleFunc("/login", handler.Login).Methods("POST", "GET")
	// s.HandleFunc("/register",handler.Register).Methods("POST")
	s.HandleFunc("/loginbygithub", handler.LoginByGitHub).Methods("POST", "GET")
	s.HandleFunc("/isvalid/{token}", handler.ValidToken).Methods("GET")

	//======================简易消息API========================================
	//更新消息状态，已读、未读
	s.HandleFunc("/auth/message/{id}", handler.UpdateMessage).Methods("PUT")
	//find 用户所有消息
	s.HandleFunc("/auth/message/{userName}", handler.FindAllMessage).Methods("GET")
	//未读消息总数
	s.HandleFunc("/auth/message/{userName}/count", handler.FindAllMessageCount).Methods("GET")

	//======================回复API============================================
	//find一个文章的所有评论
	s.HandleFunc("/comments/{articleId}", handler.FindAllCommentByArticleId).Methods("GET")
	//find一个评论下的所有回复
	s.HandleFunc("/replys/{commentId}", handler.FindAllReplyByCommentId).Methods("GET")

	//给一个文章发表评论
	s.HandleFunc("/auth/comment/{articleId}", handler.PostComment).Methods("POST")
	//回复某个评论
	s.HandleFunc("/auth/reply/{commentId}/{toUserName}", handler.PostReply).Methods("POST")

	//======================文章API============================================
	//更新文章
	s.HandleFunc("/auth/article/{articleId}", handler.UpdateArticle).Methods("PUT")
	//发表文章
	s.HandleFunc("/auth/article/{userId}", handler.SaveNewArticle).Methods("POST")
	//首页加载最新发表的文章
	s.HandleFunc("/article/findNewArticles", handler.FindNewArticles).Methods("GET")

	//加载用户所有文章，包括已审核和未审核
	s.HandleFunc("/auth/article/{userId}/all", handler.FindAllArticlesByUser).Methods("GET")
	s.HandleFunc("/article/{userId}/public", handler.FindAllArticlesByUser).Methods("GET")

	//===================其他API=================================================

	s.HandleFunc("/hotAuthor", handler.HotAuthor).Methods("GET")
	s.HandleFunc("/hotArticle", handler.HotArticle).Methods("GET")

	s.HandleFunc("/ping", ping).Methods("GET")

	//====================用户主页===============================================
	s1 := r.PathPrefix("").Subrouter()
	//用户主页模板
	s1.HandleFunc("/user/{userName}", handler.LoadUserHomePage).Methods("GET")
	// s1.HandleFunc("/user/{userName}", handler.UserHome).Methods("GET")

	//阅读某个用户的文章，加载公共文章模板
	s1.HandleFunc("/user/{userName}/article/{articleId}", handler.LoadArticlePage).Methods("GET")
	// s1.HandleFunc("/user/{userName}/article/{articleId}", handler.FindOneArticle).Methods("GET")

	//加载编辑、发帖模板
	s1.HandleFunc("/article/add", handler.LoadArticleAddPage).Methods("GET")
	s1.HandleFunc("/article/edit/{userName}/{articleId}", handler.LoadArticleEditPage).Methods("GET")
	// s1.HandleFunc("/article/edit/{userName}/{articleId}", handler.LoadEditArticleTemplate).Methods("GET")
	return r
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte("pong"))
}
