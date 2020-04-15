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

type ArticleRepositoty struct {
}

var ArticleRepos *ArticleRepositoty = &ArticleRepositoty{}

func (as *ArticleRepositoty) SaveMessage(msg model.Message) string {
	ds := DS.DataStore()
	defer ds.S.Close()

	col := ds.S.DB("edgex-club").C("message")
	timestap := time.Now().UnixNano() / 1000000
	msg.Created = timestap
	msg.Id = bson.NewObjectId()

	err := col.Insert(msg)
	if err != nil {
		log.Println("Insert article failed !")
	}
	return msg.Id.Hex()
}

func (as *ArticleRepositoty) UpdateMessage(id string) {
	ds := DS.DataStore()
	defer ds.S.Close()

	col := ds.S.DB("edgex-club").C("message")
	col.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"isRead": true}})
}

func (as *ArticleRepositoty) FindAllMessage(userName string) []model.Message {
	ds := DS.DataStore()
	defer ds.S.Close()
	msgs := make([]model.Message, 0)
	col := ds.S.DB("edgex-club").C("message")
	col.Find(bson.M{"toUserName": userName}).Sort("-created").Limit(50).All(&msgs)

	return msgs
}

func (as *ArticleRepositoty) FindAllMessageCount(userName string) int {
	ds := DS.DataStore()
	defer ds.S.Close()

	col := ds.S.DB("edgex-club").C("message")
	count, _ := col.Find(bson.M{"toUserName": userName, "isRead": false}).Count()

	return count
}

func (as *ArticleRepositoty) FindArticleIdByCommentId(commentId string) string {
	ds := DS.DataStore()
	defer ds.S.Close()
	var c model.Comment
	col := ds.S.DB("edgex-club").C("comment")
	col.Find(bson.M{"_id": bson.ObjectIdHex(commentId)}).One(&c)
	return c.ArticleId
}

func (as *ArticleRepositoty) FindAllReplyByCommentId(commentId string) []model.Reply {
	ds := DS.DataStore()
	defer ds.S.Close()
	replys := make([]model.Reply, 0)
	col := ds.S.DB("edgex-club").C("reply")
	col.Find(bson.M{"commentId": commentId}).Sort("created").All(&replys)
	return replys
}

func (as *ArticleRepositoty) FindAllCommentByArticleId(articleId string) []model.Comment {
	ds := DS.DataStore()
	defer ds.S.Close()
	comments := make([]model.Comment, 0)
	col := ds.S.DB("edgex-club").C("comment")
	col.Find(bson.M{"articleId": articleId}).Sort("-created").All(&comments)
	return comments
}

func (as *ArticleRepositoty) PostReply(r model.Reply) model.Reply {
	ds := DS.DataStore()
	defer ds.S.Close()
	col := ds.S.DB("edgex-club").C("reply")

	timestap := time.Now().UnixNano() / 1000000
	r.Created = timestap
	r.Id = bson.NewObjectId()

	err := col.Insert(r)
	if err != nil {
		log.Println("Insert reply failed !")
	}
	return r
}

func (as *ArticleRepositoty) PostComment(c model.Comment) model.Comment {
	ds := DS.DataStore()
	defer ds.S.Close()
	col := ds.S.DB("edgex-club").C("comment")

	timestap := time.Now().UnixNano() / 1000000
	c.Created = timestap
	c.Id = bson.NewObjectId()

	err := col.Insert(c)
	if err != nil {
		log.Println("Insert comment failed !")
	}
	return c
}

func (as *ArticleRepositoty) FindNewArticles() []model.Article {
	ds := DS.DataStore()
	defer ds.S.Close()
	articles := make([]model.Article, 0)
	col := ds.S.DB("edgex-club").C("article")
	col.Find(bson.M{"approved": true}).Sort("-created").Limit(50).All(&articles)
	return articles
}

func (as *ArticleRepositoty) FindAllArticles() []model.Article {
	ds := DS.DataStore()
	defer ds.S.Close()
	articles := make([]model.Article, 0)
	col := ds.S.DB("edgex-club").C("article")
	col.Find(bson.M{"approved": true}).Sort("-created").Limit(50).All(&articles)
	return articles
}

func (as *ArticleRepositoty) UpdateOne(id string, a model.Article) {
	ds := DS.DataStore()
	defer ds.S.Close()

	col := ds.S.DB("edgex-club").C("article")
	timestap := time.Now().UnixNano() / 1000000

	//bson.M{"$set": bson.M{"name": n, "modified": makeTimestamp()}
	err := col.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"title": a.Title, "content": a.Content, "type": a.Type, "posted": a.Posted, "approved": a.Approved, "modified": timestap}})
	if err != nil {
		log.Println("update article failed !")
		log.Println(err)
	}
}

func (as *ArticleRepositoty) Exists(articleId string, userId string) bool {
	ds := DS.DataStore()
	defer ds.S.Close()

	col := ds.S.DB("edgex-club").C("article")
	count, err := col.Find(bson.M{"_id": bson.ObjectIdHex(articleId), "userId": userId}).Count()

	log.Println("count===")
	log.Println(count)
	if err != nil {
		log.Println("exist  failed !")

	}
	if count > 0 {
		return true
	}
	return false
}

func (as *ArticleRepositoty) SaveOne(a model.Article) string {
	ds := DS.DataStore()
	defer ds.S.Close()

	col := ds.S.DB("edgex-club").C("article")
	timestap := time.Now().UnixNano() / 1000000
	a.Created = timestap
	a.Modified = timestap
	a.Id = bson.NewObjectId()

	err := col.Insert(a)
	if err != nil {
		log.Println("Insert article failed !")
	}
	return a.Id.Hex()
}

func (as *ArticleRepositoty) FindOne(userName string, articleId string) model.Article {
	ds := DS.DataStore()
	defer ds.S.Close()
	var a model.Article
	col := ds.S.DB("edgex-club").C("article")
	col.Find(bson.M{"_id": bson.ObjectIdHex(articleId), "userName": userName}).One(&a)
	a.ReadCount = a.ReadCount + 1
	col.UpdateId(a.Id, bson.M{"$set": bson.M{"readCount": a.ReadCount}})
	return a
}

func (as *ArticleRepositoty) FindArticleUserByArticleId(articleId string) string {
	ds := DS.DataStore()
	defer ds.S.Close()
	var a model.Article
	col := ds.S.DB("edgex-club").C("article")
	col.Find(bson.M{"_id": bson.ObjectIdHex(articleId)}).One(&a)

	return a.UserName
}

func (as *ArticleRepositoty) FindAllArticlesByUser(userId string, isFilter bool) []model.Article {
	ds := DS.DataStore()
	defer ds.S.Close()
	articles := make([]model.Article, 0)
	col := ds.S.DB("edgex-club").C("article")
	if isFilter {
		col.Find(bson.M{"userId": userId, "approved": true}).All(&articles)
	} else {
		col.Find(bson.M{"userId": userId}).All(&articles)
	}

	return articles
}

func (as *ArticleRepositoty) UserArticleCount(userName string) int {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB("edgex-club").C("article")
	count, _ := coll.Find(bson.M{"userName": userName, "approved": true}).Count()

	return count
}

func (as *ArticleRepositoty) HotAuthor() (users []model.User) {
	ds := DS.DataStore()
	defer ds.S.Close()
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

	coll := ds.S.DB("edgex-club").C("article")

	pipe := coll.Pipe(operations)
	err := pipe.AllowDiskUse().All(&result)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	users = make([]model.User, 0)
	for _, v := range result {
		u := make([]model.User, 0)
		byteArr, err := json.Marshal(v["userList"])
		if err != nil {
			log.Println(err.Error())
			return users
		}
		err = json.Unmarshal(byteArr, &u)
		if err != nil {
			log.Println(err.Error())
			return users
		}
		users = append(users, u[0])
	}

	return users
}

func (as *ArticleRepositoty) HotArticle() (articles []model.Article) {
	ds := DS.DataStore()
	defer ds.S.Close()
	coll := ds.S.DB("edgex-club").C("article")
	coll.Find(bson.M{"approved": true}).Sort("-readCount").Limit(10).All(&articles)
	return articles
}
