package storage

import (
	. "../models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	FieldID       = "_id"
	FieldMail     = "mail"
	FieldPassword = "password"
	FieldStatus   = "status"
	FieldModified = "modified"
)

func FindAllUsers() ([]*User, error) {
	client, err := NewMongoClient(Db_name)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = client.Disconnect()
	}()
	return client.FindAllUsers(Collection_users)
}

func CrateUser(data *User) error {
	if data == nil {
		return errors.Errorf("No data to create user")
	}

	client, err := NewMongoClient(Db_name)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Disconnect()
	}()

	objectID, _ := primitive.ObjectIDFromHex(data.Id)

	_, err = client.FindUser(objectID, data.Mail, Collection_users)
	if err == nil {
		return errors.New("User with this mail exist !")
	}

	data.Created = time.Now()
	data.Modified = time.Now()

	return client.InsertUser(data, Collection_users)
}

func FindOneUser(id primitive.ObjectID) (*User, error) {
	client, err := NewMongoClient(Db_name)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = client.Disconnect()
	}()
	return client.FindUser(id, "", Collection_users)
}

func UpdateUser(id primitive.ObjectID, u *User) error {
	client, err := NewMongoClient(Db_name)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Disconnect()
	}()

	_, err = client.FindUser(id, "", Collection_users)
	if err != nil {
		return errors.New("User with this id not exist !")
	}

	return client.UpdateUser(id, u, Collection_users)
}

func DeleteOneUser(id primitive.ObjectID) error {
	client, err := NewMongoClient(Db_name)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Disconnect()
	}()
	return client.DeleteUser(id, Collection_users)
}
