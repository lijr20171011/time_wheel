package timewheel

import "sync"

// secondWheel ...
// 时间轮-每秒对应的任务
type secondWheel struct {
	hour    int
	minute  int
	second  int
	taskMap map[string]*taskInfo
	mu      *sync.Mutex
}

// newSecondWheel ...
// 新建时间轮-秒
func newSecondWheel(hour, minute, second int) *secondWheel {
	return &secondWheel{
		hour:    hour,
		minute:  minute,
		second:  second,
		taskMap: map[string]*taskInfo{},
		mu:      new(sync.Mutex),
	}
}

// getTask ...
// 获取任务
func (sw *secondWheel) getTask(taskName string) (*taskInfo, bool) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	if task, ok := sw.taskMap[taskName]; ok {
		return task, true
	}
	return nil, false
}

// removeTask ...
// 删除任务
func (sw *secondWheel) removeTask(taskName string) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	if _, ok := sw.taskMap[taskName]; ok {
		delete(sw.taskMap, taskName)
		return true
	}
	return false
}

// saveTask ...
// 保存任务至秒级map
func (sw *secondWheel) saveTask(tsk *taskInfo) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.taskMap[tsk.taskName] = tsk
}
