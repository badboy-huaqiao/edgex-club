package repository

import (
	"edgex-club/internal/model"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type ReplyRepositoty interface {
	Create(r model.Reply) (model.Reply, error)
	FindAllReplyByCommentId(commentId string) (replys []model.Reply, err error)
}
type defaultReplyRepositoty struct{}

func ReplyRepositotyClient() ReplyRepositoty {
	return &defaultReplyRepositoty{}
}

func (*defaultReplyRepositoty) FindAllReplyByCommentId(commentId string) (replys []model.Reply, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	collection := ds.S.DB(EdgeXClubDB).C(ReplyTable)
	if err = collection.Find(bson.M{"commentId": commentId}).Sort("-created").All(&replys); err != nil {
		return nil, err
	}
	return replys, nil
}

func (*defaultReplyRepositoty) Create(r model.Reply) (model.Reply, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	collection := ds.S.DB(EdgeXClubDB).C(ReplyTable)

	timestap := time.Now().UnixNano() / 1000000
	r.Created = timestap
	r.Id = bson.NewObjectId()

	if err := collection.Insert(r); err != nil {
		log.Println("Insert reply failed !")
		return r, err
	}

	return r, nil
}
