package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct{
	err error
	output []byte
}

func main(){
	// 执行一个cmd，让他在一个协程里去执行，让他执行两秒
	// sleep 2;echo hello;

	// 1秒的时候杀死cmd 不会有hello输出
	var(
		ctx context.Context
		cancelFunc context.CancelFunc
		cmd *exec.Cmd
		resultChan chan*result
		res *result
	)

	// 创建一个结果队列
	resultChan = make(chan*result,1000)

	ctx,cancelFunc = context.WithCancel(context.TODO())
	go func(){
		var(
			output []byte
			err error
		)
		cmd = exec.CommandContext(ctx,"C:\\cygwin64\\bin\\bash.exe","-c","sleep 2;echo hello;")

		// 执行任务，捕获输出
		output,err = cmd.CombinedOutput()

		// 任务输出结果，传给main协程
		resultChan <- &result{
			err,
			output,
		}
	}()

	// 继续向下走
	time.Sleep(1*time.Second)

	// 取消上下文 取消子协程
	cancelFunc()

	// 在main协程里等待其他协程退出 并打印消息
	res = <- resultChan

	// 打印输出
	fmt.Println(res.err,string(res.output))
}
