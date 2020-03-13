package security

import (
	. "../storage"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

const db = "web"

type AuthData struct {
	Mail     string `json:"mail" bson:"mail"`
	Password string `json:"password" bson:"password"`
}

func CheckUser(user *AuthData) (bool, error) {
	client, err := NewMongoClient(db)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = client.Disconnect()
	}()

	if len(user.Mail) == 0 || len(user.Password) == 0 {
		return false, errors.New("No log/psw data !")
	}
	filter := bson.D{{FieldMail, user.Mail}, {FieldPassword, user.Password}}

	var u *AuthData
	err = client.Db.Collection("security").FindOne(context.TODO(), filter).Decode(&u)
	if err != nil {
		return false, err
	}

	return u != nil, nil
}
