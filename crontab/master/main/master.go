package main

import (
	"flag"
	"fmt"
	"golang/crontab/master"
	"runtime"
	"time"
)

var(
	confFile string // 配置文件的路径
)

// 解析命令行参数
func initArgs(){
	// master -config ./master.json
	// params 存到哪，指定参数，参数默认值，-h输出提示
	// master -h
	flag.StringVar(&confFile,"config","./master.json","指定master.json")
	flag.Parse()
}

// 初始化线程数量
func initEvn(){
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main()  {
	var(
		err error
	)
	// 初始化命令行参数
	initArgs()

	// 初始化线程
	initEvn()

	// 加载配置
	if err = master.InitConfig(confFile); err != nil{
		goto ERR
	}

	// 任务管理器
	if err = master.InitJobMgr(); err != nil{
		goto ERR
	}

	// 启动Api HTTP服务
	if err = master.InitApiServer();err != nil{
		goto ERR;
	}

	// 正常退出
	for {
		time.Sleep(1*time.Second)
	}

	return

	ERR:
		fmt.Println(err)
}