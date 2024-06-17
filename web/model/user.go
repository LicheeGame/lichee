package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"web/cache"
	"web/config"
	"web/dao"
	"web/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID        primitive.ObjectID `json:"uid,omitempty" bson:"_id,omitempty"`
	Openid    string             `json:"-" bson:"openid"`
	NickName  string             `json:"nickName,omitempty" bson:"nickName,omitempty"`
	AvatarUrl string             `json:"avatarUrl,omitempty" bson:"avatarUrl,omitempty"`
	Province  int                `json:"province,omitempty" bson:"province,omitempty"`
	Score     int                `json:"score,omitempty" bson:"score,omitempty"`
	Token     string             `json:"token,omitempty" bson:"-"`
}

func (User) TableName() string {
	return "users"
}

// 通过openid查找用户信息
func GetUserByOpenid(appid string, openid string) (User, error) {
	coll := dao.GetDB(appid).Collection("users")
	filter := bson.D{{"openid", openid}}
	var result User
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			//没找到
			logger.Info("not find by appid")
		}
	}
	logger.Info("GetUserByOpenid: appid:%s openid:%s %v", appid, openid, result)
	return result, err
}

// 通过uid获取用户信息
func GetUserByUID(appid string, uid string) (User, error) {
	coll := dao.GetDB(appid).Collection("users")
	id, _ := primitive.ObjectIDFromHex(uid)
	filter := bson.D{{"_id", id}}
	var result User
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			//没找到
			logger.Info("not find by uid")
		}
	}
	logger.Info("GetUserByUID: appid:%s openid:%s %v", appid, uid, result)
	return result, err
}

// 增
func AddUser(appid string, openid string) (User, error) {
	coll := dao.GetDB(appid).Collection("users")
	user := User{Openid: openid}
	result, err := coll.InsertOne(context.TODO(), &user)
	if err == nil {
		user.ID = result.InsertedID.(primitive.ObjectID)
		logger.Info("AddUser: appid:%s openid:%s %v", appid, openid, user)
		return user, err
	}
	return user, nil
}

// 改
func UpdateUser(appid string, uid string, name string, url string, province int, score int) bool {
	coll := dao.GetDB(appid).Collection("users")
	id, _ := primitive.ObjectIDFromHex(uid)
	filter := bson.M{"_id": id}
	bsonMap := bson.M{}
	if name != "" {
		bsonMap["nickName"] = name
	}
	if url != "" {
		bsonMap["avatarUrl"] = url
	}
	if province != -1 {
		bsonMap["province"] = province
	}
	if score != -1 {
		bsonMap["score"] = score
	}
	update := bson.M{"$set": bsonMap}
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	logger.Info("UpdateUser appid:%s uid:%s %v", appid, uid, result)
	if err == nil && result.MatchedCount == 1 {
		if score != -1 {
			cache.SetUserScore(appid, uid, score)
		}
		return true
	}

	if err != nil {
		logger.Info(result.UpsertedID.(primitive.ObjectID).String())
	}

	return false
}

// list
func GetRankUser(appid string) ([]User, error) {

	scoreRank, err := cache.GetUserScoreRank(appid)
	if len(scoreRank) != 0 && err == nil {
		//redis有排行榜
		var results []User
		for _, scoreMember := range scoreRank {
			uer, err := GetUserByUID(appid, scoreMember.Member.(string))
			if err == nil {
				results = append(results, uer)
			}
		}
		return results, nil
	}

	//从mongodb中获取
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{"score", -1}}).SetLimit(20)
	coll := dao.GetDB(appid).Collection("users")
	cursor, err := coll.Find(context.TODO(), filter, opts)
	var results []User

	if err != nil {
		return results, err
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return results, err
	}

	for _, result := range results {
		//写redis
		cache.SetUserScore(appid, result.ID.Hex(), result.Score)
		//res, _ := bson.MarshalExtJSON(result, false, false)
		//fmt.Println(string(res))
	}
	if len(results) != 0 {
		cache.SetUserRankExpire(appid)
	}
	logger.Info("GetRankUser appid:%s  %v", appid, results)
	return results, nil
}

/*
http://localhost:8375/code2Session?appid=wxb00370e58ccf0603&code=1111
GET https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code
{
	"openid":"xxxxxx",
	"session_key":"xxxxx",
	"unionid":"xxxxx",
	"errcode":0,
	"errmsg":"xxxxx"
}
*/

/*
	token := jwtInstance.GenerateJWT(appid, 2)
	log.Println(token)
	claims := jwtInstance.ParseJWT(token)
	log.Println(claims)
*/

func Code2Session(appid string, code string) (string, error) {
	info := config.GetWechatInfo(appid)
	if info == nil {
		return "", errors.New("appid error")
	}

	resp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		info.Appid, info.Secret, code))
	if err != nil {
		return "", err //Failed to fetch session key and openId
	}

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return result["openid"], nil
}
