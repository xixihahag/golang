package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// 时间点
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
		result *mongo.InsertOneResult
		docId primitive.ObjectID
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

	// 4插入记录（bson）
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

	if result,err = collection.InsertOne(context.TODO(),record);err != nil{
		fmt.Println(err)
		return
	}

	// id:默认生成一个全局唯一ID，objectID: 12字节二进制
	docId = result.InsertedID.(primitive.ObjectID)
	fmt.Println("自增ID",docId.Hex())
}
