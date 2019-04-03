package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main()  {
	var(
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		putResp *clientv3.PutResponse
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

	if putResp,err = kv.Put(context.TODO(),"/cron/jobs/job1","bye",clientv3.WithPrevKV()); err != nil{
		fmt.Println(err)
	}else{
		fmt.Println("revision",putResp.Header.Revision)
		if putResp.PrevKv != nil{
			fmt.Println("prevvalue",string(putResp.PrevKv.Value))
		}
	}


}