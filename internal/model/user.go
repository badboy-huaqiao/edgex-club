// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package model

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name      string        `bson:"name"          json:"name"`
	Password  string        `bson:"password"      json:"password"`
	GitHubId  string        `bson:"gitHubId"      json:"gitHubId"`
	AvatarUrl string        `bson:"avatarUrl"     json:"avatarUrl"`
	Created   int64         `bson:"created"       json:"created"`
	Modified  int64         `bson:"modified"      json:"modified"`
}

//用户登录认证之后存储用户信息，去掉敏感信息，如用户名密码等
type Credentials struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatarUrl"`
}
