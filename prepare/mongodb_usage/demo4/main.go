package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

//jobName 过滤条件
type FindByJobName struct{
	jobName string `bson:"jobName"` // jobName赋值为job10
}

func main()  {
	// mongodb读出来的是bson，需要反序列化为LogRecord对象
	var(
		client *mongo.Client
		err error
		database *mongo.Database
		collection *mongo.Collection
		cond *FindByJobName
		cursor *mongo.Cursor
		record *LogRecord
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

	// 4按照jobName字段过滤，想找出jobName=job10
	cond = &FindByJobName{
		"job10",
	}

	// 5查询
	if cursor,err = collection.Find(context.TODO(),cond); err != nil{
		fmt.Println(err)
		return
	}

	// 延时释放游标
	defer cursor.Close(context.TODO())

	// 6遍历结果集
	for cursor.Next(context.TODO()){
		// 定义一个日志对象
		record = &LogRecord{}
		// 反序列化bson对象
		if err =cursor.Decode(record); err != nil{
			fmt.Println(err)
			return
		}

		fmt.Println(*record)
	}

}
