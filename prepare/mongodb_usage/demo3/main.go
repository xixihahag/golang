package main

// 插入多条记录 未完成

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime int64 `bson:"endTime"`
}

// 一条日志
type LogRecord struct{
	Jobname string `bson:"jobName"` // 任务名
	Command string `bson:"command"`// shell命令
	Err string	`bson:"err"`// 脚本错误
	Content string `bson:"content"`// 脚本输出
	TimePoint TimePoint `bson:"timePoint"`// 执行时间
}

func main()  {
	var(
		client *mongo.Client
		err error
		database *mongo.Database
		collection *mongo.Collection
		record *LogRecord
		result *mongo.InsertManyResult
		//docId []primitive.ObjectID
	)
	// 1建立连接
	if client,err = mongo.Connect(context.TODO(),options.Client().ApplyURI("mongodb://192.168.1.139:27017")); err != nil{
		fmt.Println(err)
		return
	}

	// 2选择数据库my_db
	database = client.Database("cron")

	// 3选择表my_collections
	collection = database.Collection("log")

	record = &LogRecord{
		"job10",
		"echo hello",
		"",
		"hello",
		TimePoint{
			time.Now().Unix(),
			time.Now().Unix()+10,
		},
	}

	logArr := []interface{}{record,record,record}

	if result,err = collection.InsertMany(context.TODO(),logArr);err != nil{
		fmt.Println(err)
		return
	}

	//docId := result.InsertedIDs.(primitive.ObjectID)
	result = result

}
