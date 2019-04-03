package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main(){
	var(
		expr *cronexpr.Expression
		err error
		now time.Time
		nextTime time.Time
	)

	// 0-59分钟 0-24小时 1-31天 1-12月 0-6星期

	// 0 5 10 15 ... 55 根据时间点来执行 不是根据启动程序的时间
	// 比如你03启动程序，设置每5分钟执行一次，那么下一次执行的时间是05，不是08
	// 秒 分 时 天 月份 星期 年 下面是每5秒
	if expr,err = cronexpr.Parse("*/5 * * * * * *"); err != nil{
		fmt.Println(err)
		return
	}

	// 当前时间
	now = time.Now()
	// 下次调度时间
	nextTime = expr.Next(now)

	// 等待这个定时器超时
	time.AfterFunc(nextTime.Sub(now),func(){
		fmt.Println("onTick")
	})

	time.Sleep(5*time.Second)

	expr = expr
}