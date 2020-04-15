// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package model

import "gopkg.in/mgo.v2/bson"

type Comment struct {
	Id            bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Content       string        `bson:"content"       json:"content"`
	ArticleId     string        `bson:"articleId"     json:"articleId"`
	UserName      string        `bson:"userName"      json:"userName"`
	UserAvatarUrl string        `bson:"userAvatarUrl"      json:"userAvatarUrl"`
	Created       int64         `bson:"created"       json:"created"`
}
