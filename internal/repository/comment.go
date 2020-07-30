package repository

import (
	"edgex-club/internal/model"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type CommentRepositoty interface {
	Create(c model.Comment) (model.Comment, error)
	FindArticleIdByCommentId(commentId string) (articleId string, err error)
	FindAllCommentByArticleId(articleId string) (comments []model.Comment, err error)
}
type defaultCommentRepositoty struct{}

func CommentRepositotyClient() CommentRepositoty {
	return &defaultCommentRepositoty{}
}

func (as *defaultCommentRepositoty) Create(c model.Comment) (model.Comment, error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	collection := ds.S.DB(EdgeXClubDB).C(CommentTable)

	timestap := time.Now().UnixNano() / 1000000
	c.Created = timestap
	c.Id = bson.NewObjectId()

	if err := collection.Insert(c); err != nil {
		log.Println("Insert comment failed !")
		return c, err
	}

	return c, nil
}

func (*defaultCommentRepositoty) FindArticleIdByCommentId(commentId string) (articleId string, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	var c model.Comment
	collection := ds.S.DB(EdgeXClubDB).C(CommentTable)
	if err = collection.Find(bson.M{"_id": bson.ObjectIdHex(commentId)}).One(&c); err != nil {
		return "", err
	}
	return c.ArticleId, nil
}

func (*defaultCommentRepositoty) FindAllCommentByArticleId(articleId string) (comments []model.Comment, err error) {
	ds := DS.DataStore()
	defer ds.S.Close()
	collection := ds.S.DB(EdgeXClubDB).C(CommentTable)
	if err = collection.Find(bson.M{"articleId": articleId}).Sort("-created").All(&comments); err != nil {
		return nil, err
	}
	return comments, nil
}
