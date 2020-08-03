package repository

import (
	"edgex-club/internal/model"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type MessageRepositoty interface {
	Create(model.Message) (id string, err error)
	Update(id string) error
	FetchAllByUserName(userName string) (msgs []model.Message, err error)
	FetchPageAllByUserName(userName string, start, limit int) (msgs []model.Message, err error)
	MsgSumByUserName(userName string) (int, error)
}
type defaultMessageRepositoty struct{}

func MessageRepositotyClient() MessageRepositoty {
	return &defaultMessageRepositoty{}
}

func (*defaultMessageRepositoty) Create(msg model.Message) (id string, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()

	collection := ds.S.DB(EdgeXClubDB).C(MsgTable)
	timestap := time.Now().UnixNano() / 1000000
	msg.Created = timestap
	msg.Id = bson.NewObjectId()

	if err = collection.Insert(msg); err != nil {
		log.Println("Insert article failed !")
		return "", err
	}
	return msg.Id.Hex(), nil
}

func (*defaultMessageRepositoty) Update(id string) error {
	ds := DS.DataStore()
	defer ds.S.Close()
	collection := ds.S.DB(EdgeXClubDB).C(MsgTable)
	return collection.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"isRead": true}})
}

func (*defaultMessageRepositoty) FetchAllByUserName(userName string) (msgs []model.Message, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	collection := ds.S.DB(EdgeXClubDB).C(MsgTable)
	if err = collection.Find(bson.M{"toUserName": userName}).Sort("-created").Limit(50).All(&msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

func (*defaultMessageRepositoty) FetchPageAllByUserName(userName string, start, limit int) (msgs []model.Message, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	collection := ds.S.DB(EdgeXClubDB).C(MsgTable)
	if err = collection.Find(bson.M{"toUserName": userName}).Skip(start).Limit(limit).Sort("-created").All(&msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

func (*defaultMessageRepositoty) MsgSumByUserName(userName string) (count int, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()

	collection := ds.S.DB(EdgeXClubDB).C(MsgTable)
	return collection.Find(bson.M{"toUserName": userName, "isRead": false}).Count()
}
