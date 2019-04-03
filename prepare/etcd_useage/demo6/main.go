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
		lease clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
		putResp *clientv3.PutResponse
		kv clientv3.KV
		getResp *clientv3.GetResponse
		keepResp *clientv3.LeaseKeepAliveResponse
		keepRespChan <- chan *clientv3.LeaseKeepAliveResponse
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

	// 申请一个lease(租约)
	lease = clientv3.NewLease(client)

	// 申请一个10秒的lease
	if leaseGrantResp,err = lease.Grant(context.TODO(),10);err != nil{
		fmt.Println(err)
		return
	}

	// 自动续租 每1秒一次
	ctx,_ := context.WithTimeout(context.TODO(),5*time.Second)
	// 5秒后取消自动续租
	// 续租5秒，生命期10秒，所以一共15秒生命期
	if keepRespChan,err = lease.KeepAlive(ctx,leaseId);err != nil{
		fmt.Println(err)
		return
	}

	// 处理续租应答的协程
	go func() {
		select {
			case keepResp = <- keepRespChan:
				if keepRespChan == nil{
					fmt.Println("租约失效")
					goto END
				}else{
					fmt.Println("收到自动续租应答",keepResp.ID)
				}
		}
		END:
	}()

	// 拿到租约id
	leaseId = leaseGrantResp.ID

	// 获得kv api子集
	kv = clientv3.NewKV(client)

	// put一个kv，让它与租约关联起来，从而实现10秒后自动过期
	if putResp,err = kv.Put(context.TODO(),"/cron/lock/job1","",clientv3.WithLease(leaseId));err != nil{
		fmt.Println(err)
		return
	}

	fmt.Println("写入成功",putResp.Header.Revision)

	// 定时看下key过期了没有
	for{
		if getResp,err = kv.Get(context.TODO(),"/cron/lock/job1");err !=nil{
			fmt.Println(err)
			return
		}
		if getResp.Count == 0{
			fmt.Println("kv过期")
			break
		}

		fmt.Println("还没过期",getResp.Kvs)
		time.Sleep(2*time.Second)
	}
}
