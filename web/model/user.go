package model

import (
	"context"
	"fmt"
	"web/dao"
	"web/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"uid"`
	Openid    string             `json:"-" bson:"openid"`
	NickName  string             `json:"nickName" bson:"nickName"`
	AvatarUrl string             `json:"avatarUrl" bson:"avatarUrl"`
	Province  string             `json:"province" bson:"province"`
	Score     int                `json:"score" bson:"score"`
}

func (User) TableName() string {
	return "users"
}

func GetUserInfo(openid string) (User, error) {

	coll := dao.DB.Database("minigame").Collection("users")
	filter := bson.D{{"openid", openid}}
	var result User
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			//没找到
		}
	}
	return result, err

}

func AddUser(openid string) string {

	coll := dao.DB.Database("minigame").Collection("users")
	user := User{Openid: openid}
	result, err := coll.InsertOne(context.TODO(), &user)
	if err == nil {
		return result.InsertedID.(string)
	}
	return ""
}

func UpdateUser(openid string, name string, url string, province string, score int) {
	coll := dao.DB.Database("minigame").Collection("users")
	id, _ := primitive.ObjectIDFromHex(openid)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"nickName": name, "avatarUrl": url, "province": province, "score": score}}
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.Info(result.UpsertedID.(primitive.ObjectID).String())
	}
}

func GetRank() {
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{"score", -1}}).SetLimit(20)
	cursor, err := coll.Find(context.TODO(), filter, opts)
	var results []Course
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		res, _ := bson.MarshalExtJSON(result, false, false)
		fmt.Println(string(res))
	}
}
