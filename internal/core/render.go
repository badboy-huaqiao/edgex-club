// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package core

import (
	"html/template"
	"math"
	"strconv"
	"time"
)

const (
	TemplPathPrefix  = "static/templates/"
	BaseTemplatePath = "static/templates/base.html"
	Header_Nav       = "static/templates/header_nav.html"
)

var TemplateStore map[string]*template.Template

var dataFormate = template.FuncMap{
	"fdate": formatDate,
}

func init() {
	if TemplateStore == nil {
		TemplateStore = make(map[string]*template.Template)
	}

	indexTpl := TemplPathPrefix + "index.html"
	articleTpl := TemplPathPrefix + "/article.html"
	article_addTpl := TemplPathPrefix + "/article_add.html"
	article_editTpl := TemplPathPrefix + "/article_edit.html"
	userhomeTpl := TemplPathPrefix + "/userhome.html"

	index := template.New(indexTpl).Funcs(dataFormate)
	article := template.New(articleTpl).Funcs(dataFormate)
	article_add := template.New(article_addTpl).Funcs(dataFormate)
	article_edit := template.New(article_editTpl).Funcs(dataFormate)
	userhome := template.New(userhomeTpl).Funcs(dataFormate)

	// all, err := template.ParseGlob(TemplPathPrefix + ".html")

	TemplateStore["index"] = template.Must(index.ParseFiles(indexTpl, Header_Nav, BaseTemplatePath))
	TemplateStore["article"] = template.Must(article.ParseFiles(articleTpl, Header_Nav, BaseTemplatePath))
	TemplateStore["article_add"] = template.Must(article_add.ParseFiles(article_addTpl, Header_Nav, BaseTemplatePath))
	TemplateStore["article_edit"] = template.Must(article_edit.ParseFiles(article_editTpl, Header_Nav, BaseTemplatePath))
	TemplateStore["userhome"] = template.Must(userhome.ParseFiles(userhomeTpl, Header_Nav, BaseTemplatePath))
}

func formatDate(t int64) string {
	var byTime = []int64{365 * 24 * 60 * 60 * 1000, 12 * 24 * 60 * 60 * 1000, 24 * 60 * 60 * 1000, 60 * 60 * 1000, 60 * 1000, 1 * 1000}
	var unit = []string{"年前", "月前", "天前", "小时前", "分钟前", "秒钟前"}

	now := time.Now().UnixNano() / 1000000
	ct := now - t
	if ct < 0 {
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
