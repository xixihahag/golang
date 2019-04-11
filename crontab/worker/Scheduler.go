package worker

import (
	"fmt"
	"golang/crontab/common"
	"time"
)

// 任务调度
type Scheduler struct{
	jobEventChan chan*common.JobEvent // etcd任务事件队列
	jobPlanTable map[string]*common.JobSchedulePlan
	jobExcutingTable map[string]*common.JobExecuteInfo // 任务执行表
	jobResultChan chan *common.JobExecuteResult
}

var(
	G_scheduler *Scheduler
)

// 处理任务事件 监听etcd的事件
func (scheduler *Scheduler)handleJobEvent(jobEvent *common.JobEvent)  {
	var(
		jobSchedulePlan *common.JobSchedulePlan
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting bool
		jobExisted bool
		err error
	)
	switch jobEvent.EventType{
	case common.JOB_EVENT_SAVE:	// 保存任务事件
		if jobSchedulePlan,err = common.BuildJobSchedulePlan(jobEvent.Job);err != nil{
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case common.JOB_EVENT_DELETE:	// 删除任务事件
		if jobSchedulePlan,jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name];jobExisted{
			delete(scheduler.jobPlanTable,jobEvent.Job.Name)
		}
	case common.JOB_EVENT_KILL: // 强杀任务事件
		// 取消掉Command执行 判断任务是否在执行中
		if jobExecuteInfo,jobExecuting = scheduler.jobExcutingTable[jobEvent.Job.Name];jobExecuting{
			jobExecuteInfo.CancelFunc() // 触发command杀死shell子进程 任务退出
		}
	}
}

// 尝试执行任务
func (scheduler *Scheduler)TryStartJob(jobPlan *common.JobSchedulePlan)  {
	var(
		jobExecuteInfo *common.JobExecuteInfo
		jobExcuting bool
	)
	// 调度和执行是两件事情

	// 执行的任务可能运行很久,防止并发

	// 如果任务正在执行，那么跳过本次调度
	if jobExecuteInfo,jobExcuting = scheduler.jobExcutingTable[jobPlan.Job.Name];jobExcuting{
		return
	}

	// 构建执行状态信息
	jobExecuteInfo = common.BuildJobExcuteInfo(jobPlan)

	// 保存执行状态
	scheduler.jobExcutingTable[jobPlan.Job.Name] = jobExecuteInfo

	// 执行任务
	G_executor.ExecuteJob(jobExecuteInfo)
	fmt.Println("执行任务",jobExecuteInfo.Job.Name)
}
// 重新计算任务调度状态
func (scheduler *Scheduler)TrySchedule()(scheduleAfter time.Duration) {
	var(
		jobPlan *common.JobSchedulePlan
		now time.Time
		nearTime *time.Time
	)
	// 如果任务列表为空，随便睡眠多久
	if len(scheduler.jobPlanTable) == 0{
		scheduleAfter = 1 * time.Second
		return
	}
	// 当前时间
	now = time.Now()
	// 遍历所有任务
	for _,jobPlan = range scheduler.jobPlanTable{
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now){
			// 尝试执行任务 如果上次任务没执行完则不执行
			//fmt.Println("执行任务：",jobPlan.Job.Name)
			scheduler.TryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now)
		}

		// 统计最近一个要过期的任务事件
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime){
			nearTime = &jobPlan.NextTime
		}
	}

	// 下次调度时间 最近要执行的任务调度时间-当前时间
	scheduleAfter = (*nearTime).Sub(now)
	return
}

// 处理任务结果
func (scheduler *Scheduler)handleJobResult(result *common.JobExecuteResult)  {
	var(
		jobLog *common.JobLog
	)

	// 删除执行状态
	delete(scheduler.jobExcutingTable,result.ExcuteInfo.Job.Name)

	// 生成执行日志
	if result.Err != common.ERR_LOCK_ALREADY_REQUIRED{
		jobLog = &common.JobLog{
			JobName:result.ExcuteInfo.Job.Name,
			Command:result.ExcuteInfo.Job.Command,
			Output:string(result.Output),
			PlanTime:result.ExcuteInfo.PlanTIme.UnixNano() / 1000 / 1000,
			ScheduleTime:result.ExcuteInfo.RealTime.UnixNano() / 1000 / 1000,
			StartTime:result.StartTime.UnixNano() / 1000 / 1000,
			EndTime:result.EndTime.UnixNano() / 1000 / 1000,
		}
		if result.Err != nil{
			jobLog.Err = result.Err.Error()
		}else{
			jobLog.Err = ""
		}

		// TODO: 存储到mongodb
	}

	fmt.Println("任务执行完成",result.ExcuteInfo.Job.Name)
}

// 调度协程
func (scheduler *Scheduler)schedulerLoop()  {
	var(
		jobEvent *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
		jobResult *common.JobExecuteResult
	)
	// 初始化一次 1s
	scheduleAfter = scheduler.TrySchedule()

	// 调度的延迟定时器
	scheduleTimer = time.NewTimer(scheduleAfter)

	// 定时任务commonJob
	for{
		select {
		case jobEvent = <- scheduler.jobEventChan:	// 监听任务变化事件
			// 对内存中维护的任务列表做增删改查
			scheduler.handleJobEvent(jobEvent)
		case <- scheduleTimer.C: // 最近的任务到期了
		case jobResult = <- scheduler.jobResultChan: // 监听任务执行结果
			scheduler.handleJobResult(jobResult)
			
		}
		// 调度一次任务
		scheduleAfter = scheduler.TrySchedule()
		// 重置调度间隔
		scheduleTimer.Reset(scheduleAfter)
	}
}

// 推送任务变化事件
func (scheduler *Scheduler)PushJobEvent(jobEvent *common.JobEvent)  {
	scheduler.jobEventChan <- jobEvent
}

// 初始化调度器
func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan:make(chan *common.JobEvent,1000),
		jobPlanTable:make(map[string]*common.JobSchedulePlan),
		jobExcutingTable:make(map[string]*common.JobExecuteInfo),
		jobResultChan:make(chan *common.JobExecuteResult,1000),
	}

	// 启动调度协程
	go G_scheduler.schedulerLoop()
	return
}

// 回传任务执行结果
func (scheduler *Scheduler)PushJobResult(jobResult *common.JobExecuteResult)  {
	scheduler.jobResultChan <- jobResult
}