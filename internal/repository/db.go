// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package repository

import (
	"edgex-club/internal/config"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	EdgeXClubDB  = "edgex-club"
	ArticleTable = "article"
	MsgTable     = "message"
	CommentTable = "comment"
	ReplyTable   = "reply"
	UserTable    = "user"
)

type DataStore struct {
	S *mgo.Session
}

var DS *DataStore = &DataStore{}

func (ds *DataStore) DataStore() *DataStore {
	return &DataStore{ds.S.Copy()}
}

func DBConnect() error {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{fmt.Sprintf("%s:%d", config.Config.Database.Host, config.Config.Database.Port)},
		Timeout:  time.Duration(5000) * time.Millisecond,
		Database: config.Config.Database.DatabaseName,
		Username: config.Config.Database.Username,
		Password: config.Config.Database.Password,
	}
	s, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		return err
	}
	s.SetSocketTimeout(time.Duration(5000) * time.Millisecond)
	DS.S = s
	return nil
}
