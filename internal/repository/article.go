// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package repository

import (
	"edgex-club/internal/model"
	"encoding/json"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type ArticleRepositoty interface {
	Create(article model.Article) (id string, err error)
	UpdateOne(id string, a model.Article) error
	FetchAll() (articles []model.Article, err error)
	ExistsByUserId(articleId string, userId string) (bool, error)
	FindOne(userName string, articleId string) (a model.Article, err error)
	FindArticleUserByArticleId(articleId string) string
	FindAllArticlesByUser(userId string, isFilter bool) (articles []model.Article, err error)
	UserArticleCount(userName string) (count int, err error)
	HotAuthor() (users []model.User, err error)
	HotArticle() (articles []model.Article, err error)
}

type defaultArticleRepositoty struct{}

func ArticleRepositotyClient() ArticleRepositoty {
	return &defaultArticleRepositoty{}
}

// func (as *ArticleRepositoty2) FindNewArticles() []model.Article {
// 	ds := DS.DataStore()
// 	defer ds.S.Close()
// 	articles := make([]model.Article, 0)
// 	col := ds.S.DB(EdgeXClubDB).C(ArticleTable)
// 	col.Find(bson.M{"approved": true}).Sort("-created").Limit(50).All(&articles)
// 	return articles
// }

func (as *defaultArticleRepositoty) FetchAll() (articles []model.Article, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	col := ds.S.DB(EdgeXClubDB).C(ArticleTable)
	if err = col.Find(bson.M{"approved": true}).Sort("-created").Limit(50).All(&articles); err != nil {
		return nil, err
	}
	return articles, nil
}

func (as *defaultArticleRepositoty) UpdateOne(id string, a model.Article) error {
	ds := DS.DataStore()
	defer ds.S.Close()

	col := ds.S.DB(EdgeXClubDB).C(ArticleTable)
	timestap := time.Now().UnixNano() / 1000000

	//bson.M{"$set": bson.M{"name": n, "modified": makeTimestamp()}
	err := col.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"title": a.Title, "content": a.Content, "type": a.Type, "posted": a.Posted, "approved": a.Approved, "modified": timestap}})
	if err != nil {
		log.Println("update article failed !")
		return err
	}
	return nil
}

func (as *defaultArticleRepositoty) ExistsByUserId(articleId string, userId string) (bool, error) {
	ds := DS.DataStore()
	defer ds.S.Close()

	col := ds.S.DB(EdgeXClubDB).C(ArticleTable)
	count, err := col.Find(bson.M{"_id": bson.ObjectIdHex(articleId), "userId": userId}).Count()

	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (as *defaultArticleRepositoty) Create(a model.Article) (id string, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()

	col := ds.S.DB(EdgeXClubDB).C(ArticleTable)
	timestap := time.Now().UnixNano() / 1000000
	a.Created = timestap
	a.Modified = timestap
	a.Id = bson.NewObjectId()

	if err = col.Insert(a); err != nil {
		log.Println("Insert article failed !")
		return "", err
	}

	return a.Id.Hex(), nil
}

func (as *defaultArticleRepositoty) FindOne(userName string, articleId string) (a model.Article, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	col := ds.S.DB(EdgeXClubDB).C(ArticleTable)
	if err = col.Find(bson.M{"_id": bson.ObjectIdHex(articleId), "userName": userName}).One(&a); err != nil {
		return a, err
	}
	a.ReadCount = a.ReadCount + 1
	if err = col.UpdateId(a.Id, bson.M{"$set": bson.M{"readCount": a.ReadCount}}); err != nil {
		return a, err
	}
	return a, nil
}

func (as *defaultArticleRepositoty) FindArticleUserByArticleId(articleId string) string {
	ds := DS.DataStore()
	defer ds.S.Close()
	var a model.Article
	col := ds.S.DB(EdgeXClubDB).C(ArticleTable)
	col.Find(bson.M{"_id": bson.ObjectIdHex(articleId)}).One(&a)
	return a.UserName
}

func (as *defaultArticleRepositoty) FindAllArticlesByUser(userId string, isFilter bool) (articles []model.Article, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	col := ds.S.DB(EdgeXClubDB).C(ArticleTable)

	if isFilter {
		err = col.Find(bson.M{"userId": userId, "approved": true}).All(&articles)
	} else {
		err = col.Find(bson.M{"userId": userId}).All(&articles)
	}

	if err != nil {
		return articles, err
	}

	return articles, nil
}

func (as *defaultArticleRepositoty) UserArticleCount(userName string) (count int, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	col := ds.S.DB(EdgeXClubDB).C(ArticleTable)
	return col.Find(bson.M{"userName": userName, "approved": true}).Count()
}

func (as *defaultArticleRepositoty) HotAuthor() (users []model.User, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	col := ds.S.DB(EdgeXClubDB).C(ArticleTable)
	// db.article.aggregate([
	// 	{ $group: {
	// 		   _id: "$userName",
	// 		   amt: { $sum: 1 }
	// 	  }
	// 	},
	// 	{ $limit : 5 },
	// 	{ $sort : { amt : -1} },
	// 	{ $lookup: {from: "user",localField: "_id",foreignField: "name",as:"userList"}}
	//    ])
	result := make([]map[string]interface{}, 0)

	o1 := bson.M{"$group": bson.M{"_id": "$userName", "amt": bson.M{"$sum": 1}}}
	o2 := bson.M{"$limit": 10}
	o3 := bson.M{"$sort": bson.M{"amt": -1}}
	o4 := bson.M{"$lookup": bson.M{"from": "user", "localField": "_id", "foreignField": "name", "as": "userList"}}
	operations := []bson.M{o1, o2, o3, o4}

	pipe := col.Pipe(operations)
	if err = pipe.AllowDiskUse().All(&result); err != nil {
		return nil, err
	}

	users = make([]model.User, 0)
	for _, v := range result {
		u := make([]model.User, 0)
		byteArr, err := json.Marshal(v["userList"])
		if err != nil {
			return users, err
		}
		err = json.Unmarshal(byteArr, &u)
		if err != nil {
			return users, err
		}
		users = append(users, u[0])
	}

	return users, nil
}

func (as *defaultArticleRepositoty) HotArticle() (articles []model.Article, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	col := ds.S.DB(EdgeXClubDB).C(ArticleTable)
	if err = col.Find(bson.M{"approved": true}).Sort("-readCount").Limit(10).All(&articles); err != nil {
		return nil, err
	}
	return articles, nil
}
