package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

func main()  {
	var(
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		getResp *clientv3.GetResponse
		watchStartRevision int64
		watcher clientv3.Watcher
		watchRespChan <- chan clientv3.WatchResponse
		watchResp clientv3.WatchResponse
		event *clientv3.Event
	)

	config = clientv3.Config{
		Endpoints:[]string{"192.168.1.139:2379"},
		DialTimeout:5*time.Second,
	}

	// 建立一个客户端
	if client,err = clientv3.New(config); err != nil{
		fmt.Println(err)
		return
	}

	// 用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	// 模拟etcd中kv的变化
	go func() {
		for{
			kv.Put(context.TODO(),"/cron/jobs/job7","job7")

			kv.Delete(context.TODO(),"/cron/jobs/job7")

			time.Sleep(1*time.Second)
		}
	}()

	// 先get到当前的值，并监听后续变化
	if getResp,err = kv.Get(context.TODO(),"/cron/jobs/job7"); err != nil{
		fmt.Println(err)
		return
	}

	// 现在key是存在的
	if len(getResp.Kvs) != 0{
		fmt.Println("当前值",string(getResp.Kvs[0].Value))
	}

	// 当前etcd集群事务ID，单调递增
	watchStartRevision = getResp.Header.Revision + 1

	// 创建一个监听器
	watcher = clientv3.NewWatcher(client)

	fmt.Println("从该版本开始向后监听",watchStartRevision)

	// 启动监听 5秒后关闭
	ctx,cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
	})
	watchRespChan = watcher.Watch(ctx,"/cron/jobs/job7",clientv3.WithRev(watchStartRevision))

	// 处理kv变化事件
	for watchResp = range watchRespChan{
		for _,event = range watchResp.Events{
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为",string(event.Kv.Value),event.Kv.CreateRevision,event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了",event.Kv.ModRevision)

			}
		}
	}
}