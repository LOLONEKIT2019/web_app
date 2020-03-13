package models

import (
	"regexp"
	"time"
)

const (
	UserStatusSubscribe    = "SUBSCRIBE"
	UserStatusNotSubscribe = "NOTSUBSCRIBE"

	UserRoleUser  = "USER"
	UserRoleAdmin = "ADMIN"
)

type User struct {
	Id       string    `json:"id" bson:"_id,omitempty"`
	Mail     string    `json:"mail" bson:"mail"`
	Password string    `json:"password" bson:"password"`
	Status   string    `json:"status" bson:"status"`
	Created  time.Time `json:"created" bson:"created"`
	Modified time.Time `json:"modified" bson:"modified"`
}

func (u *User) CheckMail() bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(u.Mail)
}

func (u *User) CheckPassword() bool {

	return 8 < len(u.Password) && len(u.Password) < 16
}

func (u *User) CheckStatus() bool {
	return u.Status != UserStatusNotSubscribe && u.Status != UserStatusSubscribe
}
