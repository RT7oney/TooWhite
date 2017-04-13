package db

import (
	"TooWhite/conf"
	// "errors"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Group struct {
	Name    string
	Token   string
	Creater string
	Users   []string
}

type User struct {
	Name     string
	Token    string
	IsOnline int
	Groups   []string
}

type Msg struct {
	MsgType    int
	SendFrom   string
	SendTo     string
	SendTime   int
	Content    string
	SendStatus int
}

func newDB() *mgo.Session {
	session, err := mgo.Dial(conf.DB_DOMAIN + ":" + conf.DB_PORT)
	if err != nil {
		fmt.Println(err)
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session
}

func UserJoin(user *User) *User {
	session := newDB()
	defer session.Close()
	c := session.DB(conf.DB_DATABASE).C("user")
	result := User{}
	err := c.Find(bson.M{"token": user.Token}).One(&result)
	if err != nil {
		fmt.Println("UserJoin-51", err)
		user.IsOnline = 1
		err = c.Insert(user)
		if err != nil {
			fmt.Println("UserJoin-55", err)
		}
		return user
	}
	result.IsOnline = 1
	c.Update(bson.M{"token": user.Token},
		bson.M{"$set": bson.M{
			"isonline": 1,
			"name":     user.Name,
		}})
	return &result
}

func NewGroup(group *Group) {
	session := newDB()
	defer session.Close()
	c := session.DB(conf.DB_DATABASE).C("group")
	err := c.Insert(group)
	if err != nil {
		fmt.Println("NewGroup-76", err)
	}
	c = session.DB(conf.DB_DATABASE).C("user")
	err = c.Update(bson.M{"token": group.Creater},
		bson.M{"$push": bson.M{
			"groups": group.Token,
		}})
	if err != nil {
		fmt.Println("NewGroup-84", err)
	}
}

func UserJoinGroup(user *User, group *Group) {
	session := newDB()
	defer session.Close()
	c := session.DB(conf.DB_DATABASE).C("group")
	err := c.Update(bson.M{"token": group.Token},
		bson.M{"$push": bson.M{
			"users": user.Token,
		}})
	if err != nil {
		fmt.Println("GroupJoin-97", err)
	}
	c = session.DB(conf.DB_DATABASE).C("user")
	err = c.Update(bson.M{"token": user.Token},
		bson.M{"$push": bson.M{
			"groups": group.Token,
		}})
	if err != nil {
		fmt.Println("GroupJoin-105", err)
	}
}

func GetUserByToken(token string) *User {
	session := newDB()
	defer session.Close()
	c := session.DB(conf.DB_DATABASE).C("user")
	result := User{}
	err := c.Find(bson.M{"token": token}).One(&result)
	if err != nil {
		fmt.Println("GetUserByToken-113", err)
	}
	return &result
}

func GetGroupByToken(token string) *Group {
	session := newDB()
	defer session.Close()
	c := session.DB(conf.DB_DATABASE).C("group")
	result := Group{}
	err := c.Find(bson.M{"token": token}).One(&result)
	if err != nil {
		fmt.Println("GetGroupByToken-125", err)
	}
	return &result
}

func GetGroupsByUserToken(token string) []string {
	session := newDB()
	defer session.Close()
	c := session.DB(conf.DB_DATABASE).C("user")
	result := User{}
	err := c.Find(bson.M{"token": token}).One(&result)
	if err != nil {
		fmt.Println("GetGroupsByUserToken-116", err)
	}
	return result.Groups
}

func GetUsersByGroupToken(token string) []string {
	session := newDB()
	defer session.Close()
	c := session.DB(conf.DB_DATABASE).C("group")
	result := Group{}
	err := c.Find(bson.M{"token": token}).One(&result)
	if err != nil {
		fmt.Println("GetUsersByGroupToken-128", err)
	}
	return result.Users
}

func UserOffLine(user *User) {
	session := newDB()
	defer session.Close()
	c := session.DB(conf.DB_DATABASE).C("user")
	c.Update(bson.M{"token": user.Token},
		bson.M{"$set": bson.M{
			"isonline": 0,
		}})
}
