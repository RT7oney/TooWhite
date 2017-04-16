package db

import (
	"TooWhite/conf"
	// "errors"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
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

type OffLineMsg struct {
	MsgType  int
	SendFrom string
	SendTo   string
	SendTime time.Time
	Content  string
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
	if !IsUserInGroup(user, group) {
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
	} else {
		fmt.Println("GroupJoin-107", "用户已经存在了")
	}

}

func IsUserInGroup(user *User, group *Group) bool {
	session := newDB()
	defer session.Close()
	c := session.DB(conf.DB_DATABASE).C("group")
	result := Group{}
	c.Find(bson.M{"token": group.Token}).One(&result)
	for _, user_token := range result.Users {
		if user.Token == user_token {
			return true
		}
	}
	return false
}

func IsUserOnline(token string) bool {
	session := newDB()
	defer session.Close()
	c := session.DB(conf.DB_DATABASE).C("user")
	result := User{}
	c.Find(bson.M{"token": token}).One(&result)
	if result.IsOnline == 1 {
		return true
	}
	return false
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

func UserOffLine(user_token string) {
	session := newDB()
	defer session.Close()
	c := session.DB(conf.DB_DATABASE).C("user")
	c.Update(bson.M{"token": user_token},
		bson.M{"$set": bson.M{
			"isonline": 0,
		}})
}

func GetUserOffLineMsg(user_token string) []OffLineMsg {
	session := newDB()
	defer session.Close()
	c := session.DB(conf.DB_DATABASE).C("offlinemsg")
	results := []OffLineMsg{}
	c.Find(bson.M{"sendto": user_token}).All(&results)
	return results
}

func SaveUserOffLineMsg(msg *OffLineMsg) {
	session := newDB()
	defer session.Close()
	c := session.DB(conf.DB_DATABASE).C("offlinemsg")
	c.Insert(msg)
}

func DelUserOffLineMsg(user_token string) {
	session := newDB()
	defer session.Close()
	c := session.DB(conf.DB_DATABASE).C("offlinemsg")
	c.Remove(bson.M{"sendto": user_token})
}
