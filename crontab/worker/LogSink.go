package worker

import (
	"context"
	"fmt"
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
	autoCommitChan chan*common.LogBatch
}

var(
	G_logSink *LogSink
)

// 写入日志
func (logSink *LogSink)saveLogs(batch *common.LogBatch)  {
	logSink.logCollection.InsertMany(context.TODO(),batch.Logs)
}

func (logSink *LogSink)writeLoop()  {
	var(
		log *common.JobLog
		logBatch *common.LogBatch // 当前批次
		commitTimer *time.Timer
		timeoutBatch *common.LogBatch // 超时批次
	)

	for{
		select {
		case log = <- logSink.logChan:
			// 把这个log写入mongodb当中
			// 每次插入需要等待mongodb的一次请求往返
			// 耗时可能因为网络慢花费较长时间
			// 解决方案 批量存入
			if logBatch == nil{
				logBatch = &common.LogBatch{}
				// 让这个批次超时自动提交 默认1秒
				commitTimer = time.AfterFunc(time.Duration(G_config.JobLogBatchSize)*time.Millisecond, func(logBatch *common.LogBatch)func() {
					// 会启动另一个协程完成这个func，可能会引起同步问题
					// 所以 只发出超时通知，不提交
					return func() {
						logSink.autoCommitChan <- logBatch
					}
				}(logBatch))
			}

			// 把新日志追加进去
			logBatch.Logs = append(logBatch.Logs,log)

			// 如果批次满了，就立刻发送
			if len(logBatch.Logs) >= G_config.JobLogBatchSize{
				// 发送日志
				logSink.saveLogs(logBatch)
				// 清空logBatch
				logBatch = nil
				// 取消定时器
				commitTimer.Stop()
			}
			case timeoutBatch = <- logSink.autoCommitChan: // 过期的批次
				// 判断过期批次是否是当前批次
				if timeoutBatch != logBatch{
					continue // 跳过已经被提交的批次
				}
				// 把批次写入mongodb中
				logSink.saveLogs(timeoutBatch)
				logBatch = nil
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
		autoCommitChan:make(chan *common.LogBatch,1000),
	}

	// 启动一个mongodb处理协程
	go G_logSink.writeLoop()

	return
}

// 发送日志
func (logSink *LogSink)Append(jobLog *common.JobLog)  {
	if logSink == nil{
		fmt.Println("logSInk == nil")
	}

	select{
	case logSink.logChan <- jobLog:
	default:
		// 队列满了就丢弃
	}
}