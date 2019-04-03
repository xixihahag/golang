package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main()  {
	var(
		client *mongo.Client
		err error
		database *mongo.Database
		collection *mongo.Collection
	)
	// 1建立连接
	if client,err = mongo.Connect(context.TODO(),options.Client().ApplyURI("mongodb://192.168.1.139:27017")); err != nil{
		fmt.Println(err)
		return
	}

	// 2选择数据库my_db
	database = client.Database("my_db")

	// 3选择表my_collections
	collection = database.Collection("my_collection")

	collection = collection
}
