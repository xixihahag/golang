package main

import (
	"fmt"
	"os/exec"
)

func main(){
	var (
		cmd *exec.Cmd
		output []byte
		err error
	)

	// 生成cmd
	cmd = exec.Command("C:\\cygwin64\\bin\\bash.exe","-c","ls -l")

	// 执行命令，捕获了子进程的输出(通过pipe)
	if output,err= cmd.CombinedOutput(); err!=nil{
		fmt.Println(err)
		return
	}

	// 打印子进程的输出
	fmt.Println(string(output))
}