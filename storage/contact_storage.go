package storage

import (
	. "../models"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	db         = "web"
	collection = "contacts"
)

func (client *MongoClient) InsertContact(c *Contact, collection string) error {

	_, err := client.Db.Collection(collection).InsertOne(context.TODO(), c)
	if err != nil {
		return err
	}
	return nil
}

func CrateContact(c *Contact) error {
	client, err := NewMongoClient(db)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Disconnect()
	}()

	return client.InsertContact(c, collection)
}

func (client *MongoClient) FindContact(id primitive.ObjectID, collection string) (*Contact, error) {
	var filter bson.D

	if id.IsZero() {
		return nil, errors.New("id == 0")
	}
	filter = bson.D{{FieldID, id}}

	var c *Contact
	err := client.Db.Collection(collection).FindOne(context.TODO(), filter).Decode(&c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func FindContact(id primitive.ObjectID) (*Contact, error) {
	client, err := NewMongoClient(db)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = client.Disconnect()
	}()

	return client.FindContact(id, collection)
}

func (client *MongoClient) FindAllContacts(owner string, collection string) ([]*Contact, error) {
	ctx := context.TODO()
	cursor, err := client.Db.Collection(collection).Find(ctx, bson.D{{FieldOwner, owner}})
	if err != nil {
		return nil, err
	}
	var contacts []*Contact
	for cursor.Next(ctx) {
		var contact *Contact
		if err := cursor.Decode(&contact); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

func FindContacts(owner string) ([]*Contact, error) {
	client, err := NewMongoClient(db)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = client.Disconnect()
	}()

	return client.FindAllContacts(owner, collection)
}

func (client *MongoClient) DeleteContact(id primitive.ObjectID, collection string) error {

	ctx := context.Background()

	_, err := client.FindContact(id, db)
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

func DeleteContact(id primitive.ObjectID) error {
	client, err := NewMongoClient(db)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Disconnect()
	}()

	return client.DeleteContact(id, collection)
}

func (client *MongoClient) DeleteContacts(owner string, collection string) error {

	ctx := context.Background()

	filter := bson.D{{FieldOwner, owner}}
	_, err := client.Db.Collection(collection).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func DeleteContacts(owner string) error {
	client, err := NewMongoClient(db)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Disconnect()
	}()

	return client.DeleteContacts(owner, collection)
}

func (client *MongoClient) UpdateContact(id primitive.ObjectID, c *Contact, collection string) error {
	if id.IsZero() || c == nil {
		return errors.New("Bad id/update data")
	}
	ctx := context.Background()

	filter := bson.D{
		{FieldID, id},
	}
	var updateTo bson.D
	switch {
	case len(c.Name) == 0 && len(c.Phone) == 0:
		return errors.New("No data to update")
	case len(c.Name) == 0 && len(c.Phone) != 0:
		updateTo = bson.D{
			{"$set", bson.D{{FieldPhone, c.Phone}, {FieldModified, time.Now()}}},
		}
	case len(c.Name) != 0 && len(c.Phone) == 0:
		updateTo = bson.D{
			{"$set", bson.D{{FieldName, c.Name}, {FieldModified, time.Now()}}},
		}
	case len(c.Name) != 0 && len(c.Phone) != 0:
		updateTo = bson.D{
			{"$set", bson.D{{FieldName, c.Name}, {FieldPhone, c.Phone}, {FieldModified, time.Now()}}},
		}

	}

	c.Modified = time.Now()
	_, err := client.Db.Collection(collection).UpdateOne(ctx, filter, updateTo)
	if err != nil {
		return err
	}
	return nil
}

func UpdateContact(id primitive.ObjectID, c *Contact) error {
	client, err := NewMongoClient(db)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Disconnect()
	}()

	return client.UpdateContact(id, c, collection)
}
