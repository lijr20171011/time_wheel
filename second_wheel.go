package timewheel

import "sync"

// SecondWheel ...
// 时间轮-每秒对应的任务
type SecondWheel struct {
	Hour    int
	Minute  int
	Second  int
	TaskMap map[string]*TaskInfo
	Mu      *sync.Mutex
}

// newSecondWheel ...
// 新建时间轮-秒
func newSecondWheel(hour, minute, second int) *SecondWheel {
	return &SecondWheel{
		Hour:    hour,
		Minute:  minute,
		Second:  second,
		TaskMap: map[string]*TaskInfo{},
		Mu:      new(sync.Mutex),
	}
}

// getTask ...
// 获取任务
func (sw *SecondWheel) getTask(taskName string) (*TaskInfo, bool) {
	sw.Mu.Lock()
	defer sw.Mu.Unlock()
	if task, ok := sw.TaskMap[taskName]; ok {
		return task, true
	}
	return nil, false
}

// removeTask ...
// 删除任务
func (sw *SecondWheel) removeTask(taskName string) bool {
	sw.Mu.Lock()
	defer sw.Mu.Unlock()
	if _, ok := sw.TaskMap[taskName]; ok {
		delete(sw.TaskMap, taskName)
		return true
	}
	return false
}

// saveTask ...
// 保存任务至秒级map
func (sw *SecondWheel) saveTask(tsk *TaskInfo) {
	sw.Mu.Lock()
	defer sw.Mu.Unlock()
	sw.TaskMap[tsk.TaskName] = tsk
}
