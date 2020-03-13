package models

import (
	"regexp"
	"time"
)

const (
	FieldName  = "name"
	FieldPhone = "phone"
	FieldOwner = "owner"
)

type Contact struct {
	Id       string    `json:"id" bson:"_id,omitempty"`
	Owner    string    `json:"owner" bson:"owner"`
	Name     string    `json:"name" bson:"name"`
	Phone    string    `json:"phone" bson:"phone"`
	Created  time.Time `json:"created" bson:"created"`
	Modified time.Time `json:"modified" bson:"modified"`
}

func (c *Contact) CheckPhone() bool {
	return regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`).MatchString(c.Phone)

}

func (c *Contact) CheckName() bool {
	return len(c.Name) > 3 && len(c.Name) < 16
}
