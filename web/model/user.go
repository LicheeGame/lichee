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
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"uid,omitempty"`
	Openid    string             `json:"-" bson:"openid"`
	NickName  string             `json:"nickName,omitempty" bson:"nickName,omitempty"`
	AvatarUrl string             `json:"avatarUrl,omitempty" bson:"avatarUrl,omitempty"`
	Province  string             `json:"province,omitempty" bson:"province,omitempty"`
	Score     int                `json:"score,omitempty" bson:"score,omitempty"`
}

func (User) TableName() string {
	return "users"
}

// 通过openid查找用户信息
func GetUserByOpenid(openid string) (User, error) {

	coll := dao.DB.Database("minigame").Collection("users")
	filter := bson.D{{"openid", openid}}
	var result User
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			//没找到
			logger.Info("没找到")
		}
	}
	fmt.Printf("GetUserByOpenid: %v\n", result)
	return result, err
}

func GetUserByUID(uid string) (User, error) {
	coll := dao.DB.Database("minigame").Collection("users")
	id, _ := primitive.ObjectIDFromHex(uid)
	filter := bson.D{{"_id", id}}
	var result User
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			//没找到
			logger.Info("没找到")
		}
	}
	fmt.Printf("GetUserByUID: %v\n", result)
	return result, err
}

// 增
func AddUser(openid string) (User, error) {
	coll := dao.DB.Database("minigame").Collection("users")
	user := User{Openid: openid}
	result, err := coll.InsertOne(context.TODO(), &user)
	if err == nil {
		fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
		//return result.InsertedID.(primitive.ObjectID).String()
		user.ID = result.InsertedID.(primitive.ObjectID)
		return user, err
	}
	fmt.Printf("AddUser: %v\n", result)
	return user, nil
}

// 改
func UpdateUser(uid string, name string, url string, province string, score int) bool {
	coll := dao.DB.Database("minigame").Collection("users")
	id, _ := primitive.ObjectIDFromHex(uid)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"nickName": name, "avatarUrl": url, "province": province, "score": score}}
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	fmt.Printf("UpdateUser: %v\n", result)
	if err == nil && result.MatchedCount == 1 {
		return true
	}

	if err != nil {
		logger.Info(result.UpsertedID.(primitive.ObjectID).String())
	}

	return false
}

// list
func GetRankUser() []User {
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{"score", -1}}).SetLimit(20)
	coll := dao.DB.Database("minigame").Collection("users")
	cursor, err := coll.Find(context.TODO(), filter, opts)

	var results []User
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		res, _ := bson.MarshalExtJSON(result, false, false)
		fmt.Println(string(res))
	}
	fmt.Printf("GetRankUser: %v\n", results)
	return results
}
