package common

import (
	"encoding/json"
	"fmt"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

// 定时任务
type Job struct {
	Name string `json:"name"` // 任务名
	Command string	`json:"command"` // shell命令
	CronExpr string	`json:"cronExpr"`	// cron表达式
}

// 任务调度计划
type JobSchedulePlan struct {
	Job *Job	// 要调度的任务信息
	Expr *cronexpr.Expression	// 解析好的cronexpr表达式
	NextTime time.Time // 下次调度时间
}

// 任务执行状态
type JobExcuteInfo struct {
	Job *Job // 任务信息
	PlanTIme time.Time // 理论调度时间
	RealTime time.Time // 实际调度时间
}

// HTTP接口应答
type Response struct{
	Errno int `json:"errno"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"` // interface 万能容器
}

// 变化事件
type JobEvent struct {
	EventType int // SAVE,DELETE
	Job *Job
}

// 任务执行结果
type JobExecuteResult struct {
	ExcuteInfo *JobExcuteInfo	// 执行状态
	Output []byte	// 脚本输出
	Err error	// 脚本错误原因
	StartTime time.Time // 启动时间
	EndTime time.Time // 结束时间
}

// 应答方法
func BuildResponse(errno int,msg string,data interface{})(resp []byte,err error)  {
	// 1 定义一个response对象
	var(
		response Response
	)

	response.Errno = errno
	response.Msg = msg
	response.Data = data

	// 2 序列化
	resp,err = json.Marshal(response)

	return
}

// 反序列化Job
func UnpackJob(value []byte)(ret *Job,err error) {
	var(
		job *Job
	)

	job = &Job{}
	if err = json.Unmarshal(value,job);err != nil{
		return
	}
	ret = job

	return
}

// 从etcd的key中提取任务名
// /cron/jobs/job10 提取出 job10
func ExtractJobName(jobKey string)(string)  {
	return strings.TrimPrefix(jobKey,JOB_SAVE_DIR)
}

//  任务变化事件有两种，更新/删除
func BuildJobEvent(eventType int, job *Job)(jobEvent *JobEvent)  {
	return &JobEvent{EventType:eventType,Job:job}
}

// 构造任务执行计划
func BuildJobSchedulePlan(job *Job)(jobSchedulePlan *JobSchedulePlan,err error)  {
	var(
		expr *cronexpr.Expression
	)

	// 解析表达式
	if job == nil{
		fmt.Println("job == nil")
	}
	if expr,err = cronexpr.Parse(job.CronExpr);err != nil{
		return
	}

	// 生成调度计划
	jobSchedulePlan = &JobSchedulePlan{
		job,
		expr,
		expr.Next(time.Now()),
	}
	return
}

func BuildJobExcuteInfo(jobSchedulePlan *JobSchedulePlan) (jobExcuteInfo *JobExcuteInfo) {
	jobExcuteInfo = &JobExcuteInfo{
		Job:jobSchedulePlan.Job,
		PlanTIme:jobSchedulePlan.NextTime,	// 计划调度时间
		RealTime:time.Now(),
	}
	return
}