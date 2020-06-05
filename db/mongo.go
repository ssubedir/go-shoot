package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Username string
	Password string
	Cluster  string
}

type MongoClient struct {
	MONGO *mongo.Client
}

func (m Mongo) Connect() *MongoClient {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://"+m.Username+":"+m.Password+"@"+m.Cluster+"/test?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal("Error connecting to MongoDB")
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Error pinging MongoDB")
	} else {
		//log.Println("Connected to MongoDB")
	}

	return &MongoClient{
		MONGO: client,
	}

}

func (m *MongoClient) Disconnect() {
	err := m.MONGO.Disconnect(context.TODO())
	if err != nil {
		log.Fatal("Error closing MONGO DB")
	} else {
		//log.Println("Disconnected from MongoDB")

	}
}
