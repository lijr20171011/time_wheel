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
type TaskInfo struct {
	Hour     int           // 任务所在时
	Minute   int           // 任务所在分
	Second   int           // 任务所在秒
	TaskName string        // 任务名称
	ExecCnt  int           // 执行次数
	Interval time.Duration // 任务间隔
	TaskType TaskType      // 任务类型
	CallFunc TaskFunc      // 任务函数
	CallArgs []interface{} // 参数
}

// 任务调用方法
type TaskFunc func(args ...interface{}) error

// NewTask ...
// 新建定时任务
func NewTask(tName string, tType TaskType, execCnt int, interval time.Duration, callFunc TaskFunc, callArgs ...interface{}) (*TaskInfo, error) {
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
	return &TaskInfo{
		TaskName: tName,
		ExecCnt:  execCnt,
		Interval: interval,
		TaskType: tType,
		CallFunc: callFunc,
		CallArgs: callArgs,
	}, nil
}
