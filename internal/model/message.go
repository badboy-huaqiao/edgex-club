// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package model

import "gopkg.in/mgo.v2/bson"

type Message struct {
	Id              bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Content         string        `bson:"content"       json:"content"`
	ArticleId       string        `bson:"articleId"     json:"articleId"`
	ArticleUserName string        `bson:"articleUserName"     json:"articleUserName"`
	ToUserName      string        `bson:"toUserName"   json:"toUserName"`
	FromUserName    string        `bson:"fromUserName" json:"fromUserName"`
	IsRead          bool          `bson:"isRead"     json:"isRead"`
	Created         int64         `bson:"created"     json:"created"`
}
