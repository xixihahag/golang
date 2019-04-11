package worker

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang/crontab/common"
	"time"
)

// mongodb存储日志
type LogSink struct {
	client *mongo.Client
	logCollection *mongo.Collection
	logChan chan*common.JobLog
}

var(
	G_logSink *LogSink
)

func (logSink *LogSink)writeLoop()  {
	var(
		log *common.JobLog
	)

	for{
		select {
		case log = <- logSink.logChan:
			// 把这个log写入mongodb当中
			//
		}
	}
}

func InitLogSink()(err error)  {
	var(
		client *mongo.Client
	)

	//client,err = mongo.Connect(context.TODO(),options.Client().ApplyURI("mongodb://192.168.1.139:27017"))

	a := time.Duration(G_config.MongodbConnectionTimeout)*time.Millisecond

	// 建立mongodb链接
	if client,err = mongo.Connect(context.TODO(),&options.ClientOptions{
		ConnectTimeout:&a,
	},options.Client().ApplyURI(G_config.MongodbUri));err != nil{
		return
	}

	// 选择db和collection
	G_logSink = &LogSink{
		client:client,
		logCollection:client.Database("cron").Collection("log"),
		logChan:make(chan *common.JobLog,1000),
	}

	// 启动一个mongodb处理协程
	go G_logSink.writeLoop()

	return
}