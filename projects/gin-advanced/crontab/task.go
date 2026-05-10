package crontab

import (
	"gin-advanced/log"
	"gin-advanced/service"
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
)

// Task 定时任务结构
type Task struct {
	Name      string // 任务名称
	Spec      string // cron表达式
	Handler   func() // 执行函数
	RunOnStart bool  // 启动时是否立即执行
}

var (
	manager      *cron.Cron
	managerOnce  sync.Once
	taskHandlers = make(map[string]func())
)

// GetManager 获取任务管理器单例
func GetManager() *cron.Cron {
	managerOnce.Do(func() {
		manager = cron.New(
			cron.WithSeconds(),                       // 支持秒级精度
			cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)), // 如果任务还在运行，跳过下次执行
			cron.WithChain(cron.Recover(cron.DefaultLogger)),            // panic 恢复
		)
	})
	return manager
}

// RegisterTask 注册定时任务
func RegisterTask(task *Task) {
	taskHandlers[task.Name] = task.Handler
}

// Start 启动所有定时任务
func Start() {
	m := GetManager()

	// 注册内置任务
	for _, task := range getBuiltInTasks() {
		_, err := m.AddFunc(task.Spec, func() {
			log.I("执行定时任务: " + task.Name)
			if handler, ok := taskHandlers[task.Name]; ok {
				safeRun(handler)
			}
		})
		if err != nil {
			log.E("注册定时任务失败: " + task.Name + ", error: " + err.Error())
			continue
		}
		if task.RunOnStart {
			go safeRun(task.Handler) // 启动时立即执行
		}
		log.I(fmt.Sprintf("注册定时任务成功: %s, spec: %s", task.Name, task.Spec))
	}

	m.Start()
	log.I("定时任务管理器已启动")
}

// Stop 停止所有定时任务
func Stop() {
	if manager != nil {
		ctx := manager.Stop()
		<-ctx.Done() // 等待所有任务完成
		log.I("定时任务管理器已停止")
	}
}

// safeRun 安全执行任务
func safeRun(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.E("定时任务panic: " + recoverToString(r))
		}
	}()
	fn()
}

func recoverToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return "unknown error"
	}
}

// getBuiltInTasks 返回内置任务列表
func getBuiltInTasks() []*Task {
	return []*Task{
		{
			Name:       "CleanAccessLogs",
			Spec:       "0 0 2 * * *", // 每天凌晨2点执行 (秒 分 时 日 月 周)
			Handler:    cleanAccessLogs,
			RunOnStart: false,
		},
		{
			Name:       "LogStatistics",
			Spec:       "0 0 */6 * * *", // 每6小时执行一次
			Handler:    logStatistics,
			RunOnStart: false,
		},
	}
}

// cleanAccessLogs 清理访问日志
func cleanAccessLogs() {
	log.I("开始清理访问日志...")
	if err := service.DeleteAccessLogsByDays(7); err != nil {
		log.E("清理访问日志失败: " + err.Error())
	} else {
		log.I("清理访问日志完成")
	}
}

// logStatistics 统计日志
func logStatistics() {
	log.I("开始统计日志...")
	// 可以在这里添加日志统计逻辑
	log.I("日志统计完成")
}
