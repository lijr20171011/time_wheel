package timewheel

import "sync"

// HourWheel ...
// 时间轮-时
type HourWheel struct {
	hour      int
	minuteMap map[int]*MinuteWheel
	mu        *sync.Mutex
}

// newHourWheel ...
// 新建时间轮-时
func newHourWheel(hour int) *HourWheel {
	hw := &HourWheel{
		hour:      hour,
		minuteMap: map[int]*MinuteWheel{},
		mu:        new(sync.Mutex),
	}
	for minute := 0; minute < MinuteCnt; minute++ {
		hw.minuteMap[minute] = newMinuteWheel(hour, minute)
	}
	return hw
}

// getMinuteWheel ...
// 获取时间轮-分
func (hw *HourWheel) getMinuteWheel(minute int) (*MinuteWheel, bool) {
	hw.mu.Lock()
	defer hw.mu.Unlock()
	if mw, ok := hw.minuteMap[minute]; ok {
		return mw, true
	}
	return nil, false
}
