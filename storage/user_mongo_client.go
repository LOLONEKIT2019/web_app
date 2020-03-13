package storage

//package main
import (
	. "../models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (client *MongoClient) DropDBs(collection string) error {

	err := client.Db.Collection(collection).Drop(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

func (client *MongoClient) InsertUser(value *User, collection string) error {

	_, err := client.Db.Collection(collection).InsertOne(context.TODO(), value)
	if err != nil {
		return err
	}
	return nil
}

func (client *MongoClient) DeleteUser(id primitive.ObjectID, collection string) error {

	ctx := context.Background()

	_, err := client.FindUser(id, "", collection)
	if err != nil {
		return err
	}

	filter := bson.D{{FieldID, id}}
	_, err = client.Db.Collection(collection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (client *MongoClient) FindUser(id primitive.ObjectID, mail string, collection string) (*User, error) {
	var filter bson.D

	if !id.IsZero() {
		filter = bson.D{{FieldID, id}}
	} else {
		if len(mail) != 0 {
			filter = bson.D{{FieldMail, mail}}
		} else {
			return nil, errors.New("No id/mail to find user !")
		}
	}

	//filter := bson.D{{FieldID, id}}
	var user *User
	err := client.Db.Collection(collection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (client *MongoClient) UpdateUser(id primitive.ObjectID, u *User, collection string) error {
	if id.IsZero() || u == nil {
		return errors.New("Bad mail/update data")
	}
	ctx := context.Background()

	filter := bson.D{
		{FieldID, id},
	}
	var updateTo bson.D
	switch {
	case len(u.Password) == 0 && len(u.Status) == 0:
		return errors.New("No data to update")
	case len(u.Password) != 0 && len(u.Status) == 0:
		updateTo = bson.D{
			{"$set", bson.D{{FieldPassword, u.Password}, {FieldModified, time.Now()}}},
		}
	case len(u.Password) == 0 && len(u.Status) != 0:
		updateTo = bson.D{
			{"$set", bson.D{{FieldStatus, u.Status}, {FieldModified, time.Now()}}},
		}
	case len(u.Password) != 0 && len(u.Status) != 0:
		updateTo = bson.D{
			{"$set", bson.D{{FieldStatus, u.Status}, {FieldPassword, u.Password}, {FieldModified, time.Now()}}},
		}

	}

	_, err := client.Db.Collection(collection).UpdateOne(ctx, filter, updateTo)
	if err != nil {
		return err
	}
	return nil
}

func (client *MongoClient) FindAllUsers(collection string) ([]*User, error) {
	ctx := context.TODO()
	cursor, err := client.Db.Collection(collection).Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var users []*User
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
