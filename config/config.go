package config

import "os"

type Cfg struct {
	DataBase           string
	UsersCollection    string
	ContactsCollection string
	TokenKey           string
	Port               string
}

func GetConfig() *Cfg {
	return &Cfg{
		DataBase:           os.Getenv("DB_NAME"),
		UsersCollection:    os.Getenv("COLLECTION_USERS"),
		ContactsCollection: os.Getenv("COLLECTION_CONTACTS"),
		TokenKey:           os.Getenv("TOKEN_KEY"),
		Port:               os.Getenv("PORT"),
	}
}
