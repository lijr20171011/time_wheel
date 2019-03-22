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

// timeWheel ...
// 时间轮-天
type timeWheel struct {
	hourMap map[int]*HourWheel   // 按小时保存任务信息
	mu      *sync.Mutex          // 锁 HourMap
	stopCh  chan int             // 停止时间轮
	tasks   map[string]*taskInfo // 保存所有任务信息
	tMu     *sync.Mutex          // Tasks 锁
	ticker  *time.Ticker         // 定期执行到期任务
}

// NewTimeWheel ...
// 新建时间轮-天
func NewTimeWheel() *timeWheel {
	tw := &timeWheel{
		hourMap: map[int]*HourWheel{},
		mu:      new(sync.Mutex),
		stopCh:  make(chan int, 1),
		tasks:   map[string]*taskInfo{},
		tMu:     new(sync.Mutex),
	}
	for hour := 0; hour < HourCnt; hour++ {
		tw.hourMap[hour] = newHourWheel(hour)
	}
	return tw
}

// AddTask ...
// 将任务添加至时间轮
func (tw *timeWheel) AddTask(tsk *taskInfo) error {
	// 校验该任务是否已存在
	if _, ok := tw.getTask(tsk.taskName); ok {
		return errors.New("该任务已存在")
	}
	// 每次加入时,计算下次执行时间
	t := time.Now().Add(tsk.interval)
	// 计算剩余转到次数,默认为1
	tsk.circle = 1
	tsk.circle += int(tsk.interval / (24 * time.Hour)) // 计算跳过次数
	// 计算下次执行时分秒
	tsk.hour = t.Hour()
	tsk.minute = t.Minute()
	tsk.second = t.Second()
	// 添加至全局map
	tw.saveTask(tsk)
	hw, ok := tw.getHourWheel(tsk.hour)
	if !ok {
		err := errors.New("获取[HourWheel]异常")
		return err
	}
	mw, ok := hw.getMinuteWheel(tsk.minute)
	if !ok {
		err := errors.New("获取[MinuteWheel]异常")
		return err
	}
	sw, ok := mw.getSecondWheel(tsk.second)
	if !ok {
		err := errors.New("获取[SecondWheel]异常")
		return err
	}
	sw.saveTask(tsk)
	return nil
}

// Start ...
// 启动任务
func (tw *timeWheel) Start() {
	// 每秒执行
	tw.ticker = time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-tw.stopCh:
				tw.ticker.Stop()
				tw.ticker = nil
				return
			case <-tw.ticker.C:
				go tw.execTasks(time.Now())
			}
		}
	}()
}

// Stop ...
// 终止任务
func (tw *timeWheel) Stop() {
	tw.stopCh <- 1
}

// ClearTask ...
// 清除指定任务
func (tw *timeWheel) ClearTask(tName string) bool {
	tsk, ok := tw.getTask(tName)
	if !ok {
		return false
	}
	hw, ok := tw.getHourWheel(tsk.hour)
	if !ok {
		return false
	}
	mw, ok := hw.getMinuteWheel(tsk.minute)
	if !ok {
		return false
	}
	sw, ok := mw.getSecondWheel(tsk.second)
	if !ok {
		return false
	}
	if !sw.removeTask(tsk.taskName) {
		return false
	}
	tw.removeTask(tName)
	return true
}

// IsTaskExist ...
// 判断任务是否已存在
func (tw *timeWheel) IsTaskExist(tName string) bool {
	tw.tMu.Lock()
	defer tw.tMu.Unlock()
	_, ok := tw.tasks[tName]
	return ok
}

// getTask ...
// 根据tsk名称获取tsk信息
func (tw *timeWheel) getTask(tName string) (*taskInfo, bool) {
	tw.tMu.Lock()
	defer tw.tMu.Unlock()
	tsk, ok := tw.tasks[tName]
	return tsk, ok
}

// removeTask ...
// 从全局map中删除任务
func (tw *timeWheel) removeTask(tName string) {
	tw.tMu.Lock()
	defer tw.tMu.Unlock()
	delete(tw.tasks, tName)
	return
}

// execTasks ...
// 执行到期任务
func (tw *timeWheel) execTasks(t time.Time) {
	// 根据当前时间获取所有待执行任务
	tasks := tw.getExpireTasks(t)
	for _, tsk := range tasks {
		// 每次转到,减少剩余转到次数
		tsk.circle--
		if tsk.circle > 0 { // 未到执行时间
			continue
		}
		// 执行任务
		go func(tsk *taskInfo) {
			err := tsk.callFunc(tsk.callArgs...)
			if err != nil {
				fmt.Printf("调用函数异常: %v ;请求参数: %v \n", tsk.callFunc, tsk.callArgs)
			}
		}(tsk)
		// 根据任务类型处理
		if tsk.taskType == MultiTask {
			tsk.execCnt--
		}
		// 清除当前任务
		tw.ClearTask(tsk.taskName)
		// 清除已完成任务 (单次任务 或 已执行完的多次任务)
		if tsk.taskType == SingleTask || (tsk.taskType == MultiTask && tsk.execCnt == 0) {
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
func (tw *timeWheel) saveTask(tsk *taskInfo) {
	tw.tMu.Lock()
	defer tw.tMu.Unlock()
	tw.tasks[tsk.taskName] = tsk
	return
}

// getHourWheel ...
// 获取时间轮-时
func (tw *timeWheel) getHourWheel(hour int) (*HourWheel, bool) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if hw, ok := tw.hourMap[hour]; ok {
		return hw, true
	}
	return nil, false
}

// getExpireTasks ...
// 根据时间获取所有待执行任务
func (tw *timeWheel) getExpireTasks(t time.Time) []*taskInfo {
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
	sw.mu.Lock()
	defer sw.mu.Unlock()
	tasks := make([]*taskInfo, 0, len(sw.taskMap))
	for key, _ := range sw.taskMap {
		tasks = append(tasks, sw.taskMap[key])
	}
	return tasks
}
