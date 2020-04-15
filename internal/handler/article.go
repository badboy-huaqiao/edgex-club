// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package handler

import (
	"edgex-club/internal/model"
	repo "edgex-club/internal/repository"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	mux "github.com/gorilla/mux"
)

type TodoPageData struct {
	UserName        string
	AvatarUrl       string
	ArticleId       string
	ArticleTitle    string
	ArticleCount    int
	ArticleModified int64
	ReadCount       int64
	MD              string
	Type            string
}

func FindAllReplyByCommentId(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	commentId := vars["commentId"]

	replys := repo.ArticleRepos.FindAllReplyByCommentId(commentId)
	result, _ := json.Marshal(&replys)

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte(result))
}

func FindAllCommentByArticleId(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	articleId := vars["articleId"]

	comments := repo.ArticleRepos.FindAllCommentByArticleId(articleId)
	result, _ := json.Marshal(&comments)

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte(result))
}

func hanleContent(content string) (bool, string) {
	keywords := [...]string{"sb", "shabi", "傻逼", "傻逼", "沙比", "鸡巴", "jb", "shit", "法轮功", "法论大法好", "中国", "共产党", "政府", "操你妈", "你妈", "曹尼玛", "草泥马", "狗屎", "狗逼", "吃屎去吧", "垃圾", "傻x", "败类", "奶子", "日你妈", "狗杂种", "ri你妈"}

	for _, v := range keywords {
		b := strings.Contains(content, v)
		if b {
			log.Println("有垃圾语言:" + v)
			return true, content
		}
	}
	content = strings.Replace(content, "<", "", -1)
	content = strings.Replace(content, ">", "", -1)
	return false, content
}

func message(toUserName string, fromUserName string, articleId string, articleUserName string, t string) {
	var msg model.Message
	msg.ArticleUserName = articleUserName
	msg.ArticleId = articleId
	msg.ToUserName = toUserName
	msg.FromUserName = fromUserName
	msg.IsRead = false
	if t == "comment" {
		msg.Content = fromUserName + "评论了您的文章"
		repo.ArticleRepos.SaveMessage(msg)
	}
	if t == "reply" {
		msg.Content = fromUserName + "回复了您的评论"
		repo.ArticleRepos.SaveMessage(msg)
	}
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	msgId := vars["id"]

	repo.ArticleRepos.UpdateMessage(msgId)

	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte("ok"))
}

func FindAllMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	userName := vars["userName"]
	msgs := repo.ArticleRepos.FindAllMessage(userName)
	result, _ := json.Marshal(&msgs)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(result)
}

func FindAllMessageCount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	userName := vars["userName"]
	msgsCount := repo.ArticleRepos.FindAllMessageCount(userName)

	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte(strconv.Itoa(msgsCount)))
}

func PostReply(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	commentId := vars["commentId"]

	userStr := r.Header.Get("inner-user")
	var creds model.Credentials
	json.Unmarshal([]byte(userStr), &creds)

	fromUserName := creds.Name
	toUserName := vars["toUserName"]
	var reply model.Reply
	json.NewDecoder(r.Body).Decode(&reply)
	isVaild, content := hanleContent(reply.Content)
	if isVaild {
		log.Println(creds.Name + ": 回复有垃圾语言")
		http.Error(w, "非法字符", 3001)
		return
	}
	reply.Content = content
	reply.CommentId = commentId
	reply.FromUserName = fromUserName
	reply.ToUserName = toUserName

	reply = repo.ArticleRepos.PostReply(reply)
	log.Println("用户：" + fromUserName + " 回复了 " + toUserName)

	go func() {
		articleId := repo.ArticleRepos.FindArticleIdByCommentId(commentId)
		articleUserName := repo.ArticleRepos.FindArticleUserByArticleId(articleId)
		message(toUserName, fromUserName, articleId, articleUserName, "reply")
	}()

	result, _ := json.Marshal(&reply)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(result)

}

func PostComment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	articleId := vars["articleId"]

	userStr := r.Header.Get("inner-user")
	var creds model.Credentials
	json.Unmarshal([]byte(userStr), &creds)

	var c model.Comment
	json.NewDecoder(r.Body).Decode(&c)

	isVaild, content := hanleContent(c.Content)
	if isVaild {
		log.Println(creds.Name + ": 评论有垃圾语言")
		http.Error(w, "非法字符", 3001)
		return
	}
	c.Content = content
	c.UserName = creds.Name
	c.ArticleId = articleId
	c.UserAvatarUrl = creds.AvatarUrl

	c = repo.ArticleRepos.PostComment(c)
	log.Println("用户：" + creds.Name + " 发起了评论")
	go func() {
		articleUserName := repo.ArticleRepos.FindArticleUserByArticleId(articleId)
		message(articleUserName, creds.Name, articleId, articleUserName, "comment")
	}()

	result, _ := json.Marshal(&c)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(result)
}

func FindNewArticles(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	articles := repo.ArticleRepos.FindNewArticles()
	result, _ := json.Marshal(&articles)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte(result))
}

func FindAllArticlesByUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	userId := vars["userId"]

	userStr := r.Header.Get("inner-user")
	var creds model.Credentials
	json.Unmarshal([]byte(userStr), &creds)

	var articles []model.Article
	if creds.Id == userId {
		articles = repo.ArticleRepos.FindAllArticlesByUser(userId, false)
	} else {
		articles = repo.ArticleRepos.FindAllArticlesByUser(userId, true)
	}

	result, _ := json.Marshal(&articles)

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte(result))
}

func FindOneArticle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)

	userName := vars["userName"]
	articleId := vars["articleId"]
	fmt.Printf("userName: %s\n", vars["userName"])
	article := repo.ArticleRepos.FindOne(userName, articleId)
	user := repo.UserRepos.FindOneByName(userName)
	articleCount := repo.ArticleRepos.UserArticleCount(userName)
	t, _ := template.ParseFiles("static/articles/article_tpl.html")
	data := TodoPageData{
		UserName:        userName,
		AvatarUrl:       user.AvatarUrl,
		ArticleId:       articleId,
		ArticleTitle:    article.Title,
		ArticleCount:    articleCount,
		ArticleModified: article.Modified,
		ReadCount:       article.ReadCount,
		MD:              article.Content,
	}
	t.Execute(w, data)
}

func SaveNewArticle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var a model.Article
	vars := mux.Vars(r)
	userId := vars["userId"]
	err := json.NewDecoder(r.Body).Decode(&a)

	userStr := r.Header.Get("inner-user")
	var creds model.Credentials
	json.Unmarshal([]byte(userStr), &creds)

	a.UserId = userId
	a.AvatarUrl = creds.AvatarUrl
	a.Approved = false

	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	result := repo.ArticleRepos.SaveOne(a)
	log.Println("用户：" + creds.Name + " 保存了文章 ")

	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte(result))
}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var a model.Article
	vars := mux.Vars(r)
	articleId := vars["articleId"]
	err := json.NewDecoder(r.Body).Decode(&a)

	userStr := r.Header.Get("inner-user")
	var creds model.Credentials
	json.Unmarshal([]byte(userStr), &creds)

	if err != nil {
		log.Printf("%s：用户提交的文章无法解析", creds.Name)
		w.WriteHeader(http.StatusBadRequest)
	}

	a.UserId = creds.Id
	a.AvatarUrl = creds.AvatarUrl
	a.Approved = false
	isExist := repo.ArticleRepos.Exists(articleId, creds.Id)

	if !isExist {
		log.Printf("非法用户尝试修改不属于自己的文章: %s", creds.Name)
		http.Error(w, "非法用户", 3001)
		return
	}

	repo.ArticleRepos.UpdateOne(articleId, a)
	log.Println("用户：" + creds.Name + " 更新了文章 " + articleId)
}

func LoadEditArticleTemplate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	articleId := vars["articleId"]
	userName := vars["userName"]

	var a model.Article
	var data TodoPageData
	if articleId == "new" { //新文章
		articleId = ""
		data = TodoPageData{
			ArticleId: articleId,
		}
	} else { //编辑已有文章
		a = repo.ArticleRepos.FindOne(userName, articleId)

		data = TodoPageData{
			ArticleId:    articleId,
			MD:           a.Content,
			ArticleTitle: a.Title,
			Type:         a.Type,
		}
	}

	t, _ := template.ParseFiles("static/articles/edit_article.html")
	t.Execute(w, data)
}

func HotAuthor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	hotAuthor := repo.ArticleRepos.HotAuthor()
	result, _ := json.Marshal(&hotAuthor)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(result)
}

func HotArticle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	hotArticle := repo.ArticleRepos.HotArticle()
	result, _ := json.Marshal(&hotArticle)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(result)
}

func saveToFile(body string, userName string, articleID string) {
	// userName = "admin"
	// articleID = "011"
	// // body, err := ioutil.ReadAll(r.Body)
	// err = os.MkdirAll("static/articles/users/"+userName, os.ModePerm) //生成多级目录
	// if err != nil {
	// 	log.Println(err)
	// }
	//
	// f, err := os.Create("static/articles/users/" + userName + "/" + articleID)
	// if err != nil {
	// 	log.Printf("Error create file: %v", err)
	// }
	// _, err = f.Write(body)
	// if err != nil {
	// 	log.Printf("Error write file: %v", err)
	// }
	// // err = ioutil.WriteFile("static/dat1.html", body, 0644)
}
