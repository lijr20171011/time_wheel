package timewheel

import "sync"

// MinuteWheel ...
// 时间轮-分
type MinuteWheel struct {
	hour      int
	minute    int
	secondMap map[int]*secondWheel
	mu        *sync.Mutex
}

// newMinuteWheel ...
// 新建时间轮-分
func newMinuteWheel(hour, minute int) *MinuteWheel {
	mw := &MinuteWheel{
		hour:      hour,
		minute:    minute,
		secondMap: map[int]*secondWheel{},
		mu:        new(sync.Mutex),
	}
	for second := 0; second < SecondCnt; second++ {
		mw.secondMap[second] = newSecondWheel(hour, minute, second)
	}
	return mw
}

// getSecondWheel ...
// 获取时间轮-秒
func (mw *MinuteWheel) getSecondWheel(second int) (*secondWheel, bool) {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	if sw, ok := mw.secondMap[second]; ok {
		return sw, true
	}
	return nil, false
}
