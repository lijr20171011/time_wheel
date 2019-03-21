package timewheel

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	HourCnt   = 24
	MinuteCnt = 60
	SecondCnt = 60
)

// TimeWheel ...
// 时间轮-天
type TimeWheel struct {
	HourMap map[int]*HourWheel
	Mu      *sync.Mutex
	StopCh  chan int
	Tasks   map[string]*TaskInfo
	TMu     *sync.Mutex
}

// NewTimeWheel ...
// 新建时间轮-天
func NewTimeWheel() *TimeWheel {
	tw := &TimeWheel{
		HourMap: map[int]*HourWheel{},
		Mu:      new(sync.Mutex),
		StopCh:  make(chan int, 1),
		Tasks:   map[string]*TaskInfo{},
		TMu:     new(sync.Mutex),
	}
	for hour := 0; hour < HourCnt; hour++ {
		tw.HourMap[hour] = newHourWheel(hour)
	}
	return tw
}

// AddTask ...
// 将任务添加至时间轮
func (tw *TimeWheel) AddTask(tsk *TaskInfo) error {
	// 校验该任务是否已存在
	if _, ok := tw.GetTask(tsk.TaskName); ok {
		return errors.New("该任务已存在")
	}
	// 每次加入时,计算下次执行时间
	t := time.Now().Add(tsk.Interval)
	tsk.Hour = t.Hour()
	tsk.Minute = t.Minute()
	tsk.Second = t.Second()
	// 添加至全局map
	tw.saveTask(tsk)
	hw, ok := tw.getHourWheel(tsk.Hour)
	if !ok {
		err := errors.New("获取[HourWheel]异常")
		return err
	}
	mw, ok := hw.getMinuteWheel(tsk.Minute)
	if !ok {
		err := errors.New("获取[MinuteWheel]异常")
		return err
	}
	sw, ok := mw.getSecondWheel(tsk.Second)
	if !ok {
		err := errors.New("获取[SecondWheel]异常")
		return err
	}
	sw.saveTask(tsk)
	return nil
}

// Start ...
// 启动任务
func (tw *TimeWheel) Start() {
	go func() {
		for {
			// 每秒执行
			timer := time.NewTimer(time.Second)
			select {
			case <-tw.StopCh:
				return
			case <-timer.C:
				go tw.execTasks(time.Now())
			}
		}
	}()
}

// Stop ...
// 终止任务
func (tw *TimeWheel) Stop() {
	tw.StopCh <- 1
}

// ClearTask ...
// 清除指定任务
func (tw *TimeWheel) ClearTask(tName string) bool {
	tsk, ok := tw.GetTask(tName)
	if !ok {
		return false
	}
	hw, ok := tw.getHourWheel(tsk.Hour)
	if !ok {
		return false
	}
	mw, ok := hw.getMinuteWheel(tsk.Minute)
	if !ok {
		return false
	}
	sw, ok := mw.getSecondWheel(tsk.Second)
	if !ok {
		return false
	}
	if !sw.removeTask(tsk.TaskName) {
		return false
	}
	tw.removeTask(tName)
	return true
}

// GetTask ...
// 根据tsk名称获取tsk信息
func (tw *TimeWheel) GetTask(tName string) (*TaskInfo, bool) {
	tw.TMu.Lock()
	defer tw.TMu.Unlock()
	tsk, ok := tw.Tasks[tName]
	return tsk, ok
}

// removeTask ...
// 从全局map中删除任务
func (tw *TimeWheel) removeTask(tName string) {
	tw.TMu.Lock()
	defer tw.TMu.Unlock()
	delete(tw.Tasks, tName)
	return
}

// execTasks ...
// 执行到期任务
func (tw *TimeWheel) execTasks(t time.Time) {
	// 根据当前时间获取所有待执行任务
	tasks := tw.getExpireTasks(t)
	for _, tsk := range tasks {
		// 执行任务
		go func(tsk *TaskInfo) {
			err := tsk.CallFunc(tsk.CallArgs...)
			if err != nil {
				fmt.Printf("调用函数异常: %v ;请求参数: %v \n", tsk.CallFunc, tsk.CallArgs)
			}
		}(tsk)
		// 根据任务类型处理
		if tsk.TaskType == MultiTask {
			tsk.ExecCnt--
		}
		// 清除当前任务
		tw.ClearTask(tsk.TaskName)
		// 清除已完成任务 (单次任务 或 已执行完的多次任务)
		if tsk.TaskType == SingleTask || (tsk.TaskType == MultiTask && tsk.ExecCnt == 0) {
			continue
		}
		// 重新加入任务队列 (未执行完的多次任务 或 永久任务)
		err := tw.AddTask(tsk)
		if err != nil {
			panic(err)
		}
	}
}

// saveTask ...
// 将任务保存至全局map
func (tw *TimeWheel) saveTask(tsk *TaskInfo) {
	tw.TMu.Lock()
	defer tw.TMu.Unlock()
	tw.Tasks[tsk.TaskName] = tsk
	return
}

// getHourWheel ...
// 获取时间轮-时
func (tw *TimeWheel) getHourWheel(hour int) (*HourWheel, bool) {
	tw.Mu.Lock()
	defer tw.Mu.Unlock()
	if hw, ok := tw.HourMap[hour]; ok {
		return hw, true
	}
	return nil, false
}

// getExpireTasks ...
// 根据时间获取所有待执行任务
func (tw *TimeWheel) getExpireTasks(t time.Time) []*TaskInfo {
	hw, ok := tw.getHourWheel(t.Hour())
	if !ok {
		return nil
	}
	mw, ok := hw.getMinuteWheel(t.Minute())
	if !ok {
		return nil
	}
	sw, ok := mw.getSecondWheel(t.Second())
	if !ok {
		return nil
	}
	sw.Mu.Lock()
	defer sw.Mu.Unlock()
	tasks := make([]*TaskInfo, 0, len(sw.TaskMap))
	for key, _ := range sw.TaskMap {
		tasks = append(tasks, sw.TaskMap[key])
	}
	return tasks
}
