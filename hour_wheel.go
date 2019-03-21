package timewheel

import "sync"

// HourWheel ...
// 时间轮-时
type HourWheel struct {
	Hour      int
	MinuteMap map[int]*MinuteWheel
	Mu        *sync.Mutex
}

// newHourWheel ...
// 新建时间轮-时
func newHourWheel(hour int) *HourWheel {
	hw := &HourWheel{
		Hour:      hour,
		MinuteMap: map[int]*MinuteWheel{},
		Mu:        new(sync.Mutex),
	}
	for minute := 0; minute < MinuteCnt; minute++ {
		hw.MinuteMap[minute] = newMinuteWheel(hour, minute)
	}
	return hw
}

// getMinuteWheel ...
// 获取时间轮-分
func (hw *HourWheel) getMinuteWheel(minute int) (*MinuteWheel, bool) {
	hw.Mu.Lock()
	defer hw.Mu.Unlock()
	if mw, ok := hw.MinuteMap[minute]; ok {
		return mw, true
	}
	return nil, false
}
