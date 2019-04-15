package master

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"golang/crontab/common"
	"time"
)

type WorkerMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease

}

var(
	G_workerMgr *WorkerMgr
)

// 获取在线worker列表
func (workerMgr *WorkerMgr)ListWorkers()(workerArr []string,err error)  {
	var(
		getResp *clientv3.GetResponse
		kv *mvccpb.KeyValue
		workerIp string
	)
	workerArr = make([]string,0)

	if getResp,err = workerMgr.kv.Get(context.TODO(),common.JOB_WORKER_DIR,clientv3.WithPrefix());err != nil{
		return
	}

	// 解析每个节点的ip
	for _,kv = range getResp.Kvs{
		// kv.Key: /cron/workers/192.168.1.1
		workerIp = common.ExtractWorkerIp(string(kv.Key))
		workerArr = append(workerArr,workerIp)
	}

	return
}

func InitWorkerMgr()(err error)  {
	var(
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
	)
	// 初始化配置
	config = clientv3.Config{
		Endpoints: G_config.EtcdEndpoints,	// 集群地址
		DialTimeout:time.Duration(G_config.EtcdDialTimeout)*time.Millisecond,	// 链接超时
	}

	// 建立连接
	if client,err = clientv3.New(config);err != nil{
		return
	}

	// 得到KV和Lease的api子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_workerMgr = &WorkerMgr{
		client:client,
		kv:kv,
		lease:lease,
	}

	return
}