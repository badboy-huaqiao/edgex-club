// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package model

import "gopkg.in/mgo.v2/bson"

type Article struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Title     string        `bson:"title"    json:"title"`
	Content   string        `bson:"content"  json:"content"`
	Type      string        `bson:"type"     json:"type"`
	UserName  string        `bson:"userName" json:"userName"`
	UserId    string        `bson:"userId"   json:"userId"`
	AvatarUrl string        `bson:"avatarUrl"     json:"avatarUrl"`
	Posted    bool          `bson:"posted"   json:"posted"`
	Approved  bool          `bson:"approved"   json:"approved"`
	ReadCount int64         `bson:"readCount" json:"readCount"`
	Created   int64         `bson:"created"  json:"created"`
	Modified  int64         `bson:"modified" json:"modified"`
}
