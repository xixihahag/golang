package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// startTime 小于某时间
// {"$lt":timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

// {"timePoint.startTime":}
type DeleteCond struct{
	beforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

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
		delCond *DeleteCond
		delResult *mongo.DeleteResult
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

	// 4要删除开始时间早于当前时间的所有日志
	// delete({"timePoint.startTime":{"$lt":当前时间}})
	delCond = &DeleteCond{
		TimeBeforeCond{time.Now().Unix()},
	}

	// 执行删除
	if delResult,err = collection.DeleteMany(context.TODO(),delCond);err != nil{
		fmt.Println(err)
		return
	}

	fmt.Println("删除的行数",delResult.DeletedCount)
}
