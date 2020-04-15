// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package repository

import (
	"edgex-club/internal/model"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type UserRepository struct {
}

var UserRepos *UserRepository = &UserRepository{}

func (ur *UserRepository) ExistsByGitHub(u model.User) (bool, error) {
	ds := DS.DataStore()
	defer ds.S.Close()

	coll := ds.S.DB("edgex-club").C("user")
	count, err := coll.Find(bson.M{"name": u.Name, "gitHubId": u.GitHubId}).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (ur *UserRepository) Exists(u model.User) (bool, error) {
	ds := DS.DataStore()
	defer ds.S.Close()

	coll := ds.S.DB("edgex-club").C("user")
	count, err := coll.Find(bson.M{"name": u.Name, "password": u.Password}).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (ur *UserRepository) FindOneByName(name string) model.User {
	ds := DS.DataStore()
	defer ds.S.Close()
	var u model.User
	coll := ds.S.DB("edgex-club").C("user")
	coll.Find(bson.M{"name": name}).One(&u)

	return u
}

func (ur *UserRepository) Insert(u model.User) (string, error) {
	ds := DS.DataStore()
	defer ds.S.Close()

	coll := ds.S.DB("edgex-club").C("user")
	timestap := time.Now().UnixNano() / 1000000
	u.Created = timestap
	u.Id = bson.NewObjectId()
	err := coll.Insert(u)
	if err != nil {
		log.Println("Insert user failed !")
		return "", err
	}

	return u.Id.Hex(), nil
}

func (ur *UserRepository) Update(user model.User) error {
	ds := DS.DataStore()
	defer ds.S.Close()

	coll := ds.S.DB("edgex-club").C("user")
	err := coll.UpdateId(user.Id, &user)

	if err != nil {
		log.Println("Update user failed !")
		return err
	}

	return nil
}

func (ur *UserRepository) Delete(id string) error {
	ds := DS.DataStore()
	defer ds.S.Close()

	coll := ds.S.DB("edgex-club").C("user")
	err := coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Println("Delete user failed!" + err.Error())
		return err
	}
	return nil
}
