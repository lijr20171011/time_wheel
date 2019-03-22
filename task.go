package timewheel

import (
	"errors"
	"time"
)

type TaskType int

const (
	SingleTask    TaskType = iota + 1 // 单次任务
	MultiTask                         // 多次任务
	PermanentTask                     // 永久任务
)

// 任务详情
type taskInfo struct {
	hour     int           // 任务所在时
	minute   int           // 任务所在分
	second   int           // 任务所在秒
	taskName string        // 任务名称
	circle   int           // 剩余转到次数,定义第n次转到开始执行
	execCnt  int           // 执行次数
	interval time.Duration // 任务间隔
	taskType TaskType      // 任务类型
	callFunc TaskFunc      // 任务函数
	callArgs []interface{} // 参数
}

// 任务调用方法
type TaskFunc func(args ...interface{}) error

// NewTask ...
// 新建定时任务
func NewTask(tName string, tType TaskType, execCnt int, interval time.Duration, callFunc TaskFunc, callArgs ...interface{}) (*taskInfo, error) {
	if tName == "" {
		return nil, errors.New("任务名称不能为空")
	}
	if tType != SingleTask && tType != MultiTask && tType != PermanentTask {
		return nil, errors.New("未知的任务类型")
	}
	if tType == MultiTask && execCnt < 1 {
		return nil, errors.New("多次任务执行次数不能少于1")
	}
	if interval < time.Second {
		return nil, errors.New("时间间隔不能少于1s")
	}
	return &taskInfo{
		taskName: tName,
		execCnt:  execCnt,
		interval: interval,
		taskType: tType,
		callFunc: callFunc,
		callArgs: callArgs,
	}, nil
}
