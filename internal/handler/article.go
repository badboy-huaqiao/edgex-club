// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package handler

import (
	"bytes"
	"edgex-club/internal/core"
	"edgex-club/internal/model"
	repo "edgex-club/internal/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	mux "github.com/gorilla/mux"
)

const (
	ContentType     string = "Content-Type"
	ContentTypeText string = "text/plain;charset=utf-8"
	ContentTypeJSON string = "application/json;charset=utf-8"
	CredsUser       string = "CredsUser"
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

func checkContentSize(data string) bool {
	fixedSize := 1024 * 1024 * 5
	result := []byte(data)
	if len(result) > fixedSize {
		return false
	}
	return true
}

func FindAllReplyByCommentId(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	commentId := vars["commentId"]

	replys, err := repo.ReplyRepositotyClient().FindAllReplyByCommentId(commentId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var result []byte
	if result, err = json.Marshal(&replys); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(ContentType, ContentTypeJSON)
	w.Write(result)
}

func FindAllCommentByArticleId(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	articleId := vars["articleId"]

	comments, err := repo.CommentRepositotyClient().FindAllCommentByArticleId(articleId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var result []byte
	if result, err = json.Marshal(&comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(ContentType, ContentTypeJSON)
	w.Write(result)
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
	msg := model.Message{
		ArticleUserName: articleUserName,
		ArticleId:       articleId,
		ToUserName:      toUserName,
		FromUserName:    fromUserName,
		IsRead:          false,
	}

	if t == "comment" {
		msg.Content = fmt.Sprintf("%s评论了您的文章", fromUserName)
	} else if t == "reply" {
		msg.Content = fmt.Sprintf("%s回复了您的评论", fromUserName)
	}
	if _, err := repo.MessageRepositotyClient().Create(msg); err != nil {
		log.Println(err.Error())
	}
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	msgId := vars["id"]

	if err := repo.MessageRepositotyClient().Update(msgId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(ContentType, ContentTypeText)
	w.Write([]byte("ok"))
}

func FindAllMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	userName := vars["userName"]
	msgs, err := repo.MessageRepositotyClient().FetchAllByUserName(userName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(&msgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(ContentType, ContentTypeJSON)
	w.Write(result)
}

func FindAllMessageCount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	userName := vars["userName"]
	msgsCount, err := repo.MessageRepositotyClient().MsgSumByUserName(userName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(ContentType, ContentTypeText)
	w.Write([]byte(strconv.Itoa(msgsCount)))
}

func PostReply(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	commentId := vars["commentId"]

	creds := genCredsUser(r)

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

	if ok := checkContentSize(reply.Content); !ok {
		http.Error(w, "内容不能超过5M", http.StatusLengthRequired)
		return
	}

	reply.Content = content
	reply.CommentId = commentId
	reply.FromUserName = fromUserName
	reply.ToUserName = toUserName

	reply, err := repo.ReplyRepositotyClient().Create(reply)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("用户：" + fromUserName + " 回复了 " + toUserName)

	go func() {
		articleId, _ := repo.CommentRepositotyClient().FindArticleIdByCommentId(commentId)
		articleUserName := repo.ArticleRepositotyClient().FindArticleUserByArticleId(articleId)
		message(toUserName, fromUserName, articleId, articleUserName, "reply")
	}()

	var bytesBuf bytes.Buffer
	t := core.TemplateStore["replyTpl"]
	if err := t.ExecuteTemplate(&bytesBuf, "reply", reply); err != nil {
		fmt.Printf("bytesBuf err : %v\n", err.Error())
	}
	w.Header().Set(ContentType, ContentTypeText)
	w.Write([]byte(bytesBuf.String()))

}

func PostComment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	articleId := vars["articleId"]

	creds := genCredsUser(r)

	var c model.Comment
	json.NewDecoder(r.Body).Decode(&c)

	isVaild, content := hanleContent(c.Content)
	if isVaild {
		log.Println(creds.Name + ": 评论有垃圾语言")
		http.Error(w, "非法字符", 3001)
		return
	}

	if ok := checkContentSize(c.Content); !ok {
		http.Error(w, "内容不能超过5M", http.StatusLengthRequired)
		return
	}
	c.Content = content
	c.UserName = creds.Name
	c.ArticleId = articleId
	c.UserAvatarUrl = creds.AvatarUrl

	c, err := repo.CommentRepositotyClient().Create(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("用户：" + creds.Name + " 发起了评论")
	go func() {
		articleUserName := repo.ArticleRepositotyClient().FindArticleUserByArticleId(articleId)
		message(articleUserName, creds.Name, articleId, articleUserName, "comment")
	}()

	t := core.TemplateStore["commentTpl"]
	t.ExecuteTemplate(w, "comment", c)
}

func FindNewArticles(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	articles, err := repo.ArticleRepositotyClient().FetchAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result, _ := json.Marshal(&articles)
	w.Header().Set(ContentType, ContentTypeJSON)
	w.Write(result)
}

func FindAllArticlesByUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	userId := vars["userId"]

	creds := genCredsUser(r)

	var articles []model.Article
	var err error
	var filter bool

	if creds.Id != userId {
		filter = true
	}
	if articles, err = repo.ArticleRepositotyClient().FindAllArticlesByUser(userId, filter); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, _ := json.Marshal(&articles)

	w.Header().Set(ContentType, ContentTypeJSON)
	w.Write(result)
}

func SaveNewArticle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var a model.Article
	vars := mux.Vars(r)
	userId := vars["userId"]
	err := json.NewDecoder(r.Body).Decode(&a)

	if ok := checkContentSize(a.Content); !ok {
		http.Error(w, "内容不能超过5M", http.StatusLengthRequired)
		return
	}

	creds := genCredsUser(r)

	a.UserId = userId
	a.AvatarUrl = creds.AvatarUrl
	a.Approved = false

	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	id, err := repo.ArticleRepositotyClient().Create(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("用户：" + creds.Name + " 保存了文章 ")

	w.Header().Set(ContentType, ContentTypeText)
	w.Write([]byte(id))
}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var a model.Article
	vars := mux.Vars(r)
	articleId := vars["articleId"]

	creds := genCredsUser(r)

	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		log.Printf("%s：用户提交的文章无法解析", creds.Name)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ok := checkContentSize(a.Content); !ok {
		http.Error(w, "内容不能超过5M", http.StatusLengthRequired)
		return
	}

	a.UserId = creds.Id
	a.AvatarUrl = creds.AvatarUrl
	a.Approved = false
	isExist, _ := repo.ArticleRepositotyClient().ExistsByUserId(articleId, creds.Id)

	if !isExist {
		log.Printf("非法用户尝试修改不属于自己的文章: %s", creds.Name)
		http.Error(w, "非法用户", 3001)
		return
	}

	if err := repo.ArticleRepositotyClient().UpdateOne(articleId, a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("用户：" + creds.Name + " 更新了文章 " + articleId)
}

func HotAuthor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	hotAuthor, _ := repo.ArticleRepositotyClient().HotAuthor()
	result, err := json.Marshal(&hotAuthor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set(ContentType, ContentTypeJSON)
	w.Write(result)
}

func HotArticle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	hotArticle, _ := repo.ArticleRepositotyClient().HotArticle()
	result, err := json.Marshal(&hotArticle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set(ContentType, ContentTypeJSON)
	w.Write(result)
}
