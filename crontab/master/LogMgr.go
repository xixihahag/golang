package master

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang/crontab/common"
	"time"
)

// mongodb日志管理
type LogMgr struct {
	client *mongo.Client
	logCollection *mongo.Collection
}

var(
	G_logMgr *LogMgr
)

func InitLogMgr()(err error)  {
	var(
		client *mongo.Client
	)

	//client,err = mongo.Connect(context.TODO(),options.Client().ApplyURI("mongodb://192.168.1.139:27017"))

	a := time.Duration(G_config.MongodbConnectionTimeout)*time.Millisecond

	//建立mongodb链接
	if client,err = mongo.Connect(context.TODO(),&options.ClientOptions{
		ConnectTimeout:&a,
	},options.Client().ApplyURI(G_config.MongodbUri));err != nil{
		return
	}

	G_logMgr = &LogMgr{
		client:client,
		logCollection:client.Database("cron").Collection("log"),
	}
	return
}

// 查看任务日志
func (logMgr *LogMgr)ListLog(name string, skip int64,limit int64)(logArr []*common.JobLog,err error)  {
	var(
		filter *common.JobLogFilter
		logSort *common.SortLogByStartTime
		cursor *mongo.Cursor
		jobLog *common.JobLog
	)


	// len(logArr)
	logArr = make([]*common.JobLog,0)

	// 过滤条件
	filter = &common.JobLogFilter{
		JobName:name,
	}

	//fmt.Println("name =",filter.JobName)

	// 按照任务开始时间倒排
	logSort = &common.SortLogByStartTime{SortOrder:-1}

	opt := &options.FindOptions{
		Sort:logSort,
		Skip:&skip,
		Limit:&limit,
	}
	if cursor,err = logMgr.logCollection.Find(context.TODO(),filter,opt); err != nil{
		return
	}

	// 延迟释放游标
	defer cursor.Close(context.TODO())

	//fmt.Println("*")

	for cursor.Next(context.TODO()){
		//fmt.Println("&")
		jobLog = &common.JobLog{}

		// 反序列化BSON
		if err = cursor.Decode(jobLog);err != nil{
			//fmt.Println("不合法")
			continue //有日志不合法 直接忽略
		}

		//fmt.Println("add")
		logArr = append(logArr, jobLog)
	}

	return
}