// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package repository

import (
	"edgex-club/internal/model"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type UserRepository interface {
	Add(u model.User) (id string, err error)
	Update(u model.User) error
	FetchOneByGitHub(githubId string) (u model.User, err error)
	FetchOneByName(name string) (u model.User, err error)
}

func UserRepositoryClient() UserRepository {
	return &defaultUserRepository{}
}

type defaultUserRepository struct {
}

func (*defaultUserRepository) Add(u model.User) (id string, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	col := ds.S.DB(EdgeXClubDB).C(UserTable)
	timestap := time.Now().UnixNano() / 1000000
	u.Created = timestap
	u.Id = bson.NewObjectId()
	if err = col.Insert(u); err != nil {
		log.Printf("Insert user failed: %s", err.Error())
		return "", err
	}
	return u.Id.Hex(), nil
}

func (*defaultUserRepository) Update(u model.User) error {
	ds := DS.DataStore()
	defer ds.S.Close()
	ts := time.Now().UnixNano() / 1000000
	col := ds.S.DB(EdgeXClubDB).C(UserTable)
	if err := col.UpdateId(u.Id, bson.M{"$set": bson.M{"avatarUrl": u.AvatarUrl, "Modified": ts}}); err != nil {
		log.Printf("Update user failed: %s\n", err.Error())
		return err
	}
	return nil
}

func (*defaultUserRepository) FetchOneByGitHub(gitHubId string) (u model.User, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	col := ds.S.DB(EdgeXClubDB).C(UserTable)
	if err = col.Find(bson.M{"gitHubId": gitHubId}).One(&u); err != nil {
		return u, err
	}
	return u, nil
}

func (*defaultUserRepository) FetchOneByName(name string) (u model.User, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	col := ds.S.DB(EdgeXClubDB).C(UserTable)

	if err = col.Find(bson.M{"name": name}).One(&u); err != nil {
		return u, err
	}
	return u, nil
}
