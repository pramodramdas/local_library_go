package config

import (
	"context"
	"digitalLibrary/utils"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

var newCustomLog = utils.NewCustomLog

func ConnectToDB() {
	//var collection *mongo.Collection
	var ctx = context.TODO()

	clientOpts := options.Client().ApplyURI(os.Getenv("MONGO_URL"))

	clientOpts = clientOpts.SetMaxPoolSize(50)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database(os.Getenv("MONGO_DB_NAME"))
	newCustomLog("info", "db connected", make([]byte, 0), make([]byte, 0), nil)
}

func GetDb() *mongo.Database {
	return DB
}
