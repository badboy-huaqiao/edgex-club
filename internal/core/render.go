// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package core

import (
	"html/template"
	"math"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	TemplPathPrefix  = "static/templates/"
	BaseTemplatePath = "static/templates/base.html"
	Header_Nav_Path  = "static/templates/header_nav.html"
)

var TemplateStore map[string]*template.Template

var funcs = template.FuncMap{
	"fdate":     formatDate,
	"bsonIdStr": convertBsonToStr,
}

func init() {
	if TemplateStore == nil {
		TemplateStore = make(map[string]*template.Template)
	}

	indexTplPath := TemplPathPrefix + "index.html"
	articleTplPath := TemplPathPrefix + "article.html"
	articleAddTplPath := TemplPathPrefix + "article_add.html"
	articleEditTplPath := TemplPathPrefix + "article_edit.html"
	userhomeTplPath := TemplPathPrefix + "userhome.html"

	commentTplPath := TemplPathPrefix + "comment.tpl"
	replyTplPath := TemplPathPrefix + "reply.tpl"

	errorTplPath := TemplPathPrefix + "error.html"

	indexTpl := template.New(indexTplPath).Funcs(funcs)
	articleTpl := template.New(articleTplPath).Funcs(funcs)
	articleAddTpl := template.New(articleAddTplPath).Funcs(funcs)
	articleEditTpl := template.New(articleEditTplPath).Funcs(funcs)
	userhomeTpl := template.New(userhomeTplPath).Funcs(funcs)

	commentTpl := template.New(commentTplPath).Funcs(funcs)
	replyTpl := template.New(replyTplPath).Funcs(funcs)

	// all, err := template.ParseGlob(TemplPathPrefix + ".html")

	TemplateStore["index"] = template.Must(indexTpl.ParseFiles(indexTplPath, Header_Nav_Path, BaseTemplatePath))
	TemplateStore["article"] = template.Must(articleTpl.ParseFiles(articleTplPath, Header_Nav_Path, BaseTemplatePath))
	TemplateStore["article_add"] = template.Must(articleAddTpl.ParseFiles(articleAddTplPath, Header_Nav_Path, BaseTemplatePath))
	TemplateStore["article_edit"] = template.Must(articleEditTpl.ParseFiles(articleEditTplPath, Header_Nav_Path, BaseTemplatePath))
	TemplateStore["userhome"] = template.Must(userhomeTpl.ParseFiles(userhomeTplPath, Header_Nav_Path, BaseTemplatePath))

	TemplateStore["commentTpl"] = template.Must(commentTpl.ParseFiles(commentTplPath))
	TemplateStore["replyTpl"] = template.Must(replyTpl.ParseFiles(replyTplPath))
	TemplateStore["errorTpl"] = template.Must(template.ParseFiles(errorTplPath))
}

func convertBsonToStr(bsonId bson.ObjectId) string {
	return bsonId.Hex()
}

func formatDate(t int64) string {
	var byTime = []int64{365 * 24 * 60 * 60 * 1000, 12 * 24 * 60 * 60 * 1000, 24 * 60 * 60 * 1000, 60 * 60 * 1000, 60 * 1000, 1 * 1000}
	var unit = []string{"年前", "月前", "天前", "小时前", "分钟前", "秒钟前"}

	now := time.Now().UnixNano() / 1000000
	ct := now - t
	if ct <= 1000 { // 1秒以内
		return "刚刚"
	}
	var res string
	for i := 0; i < len(byTime); i++ {
		if ct < byTime[i] {
			continue
		}
		var temp = math.Floor(float64(ct / byTime[i]))
		ct = ct % byTime[i]
		if temp > 0 {
			var tempStr string
			tempStr = strconv.FormatFloat(temp, 'f', -1, 64)
			res = tempStr + unit[i]
		}
		break
	}
	return res
}
