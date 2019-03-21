package timewheel

import "sync"

// MinuteWheel ...
// 时间轮-分
type MinuteWheel struct {
	Hour      int
	Minute    int
	SecondMap map[int]*SecondWheel
	Mu        *sync.Mutex
}

// newMinuteWheel ...
// 新建时间轮-分
func newMinuteWheel(hour, minute int) *MinuteWheel {
	mw := &MinuteWheel{
		Hour:      hour,
		Minute:    minute,
		SecondMap: map[int]*SecondWheel{},
		Mu:        new(sync.Mutex),
	}
	for second := 0; second < SecondCnt; second++ {
		mw.SecondMap[second] = newSecondWheel(hour, minute, second)
	}
	return mw
}

// getSecondWheel ...
// 获取时间轮-秒
func (mw *MinuteWheel) getSecondWheel(second int) (*SecondWheel, bool) {
	mw.Mu.Lock()
	defer mw.Mu.Unlock()
	if sw, ok := mw.SecondMap[second]; ok {
		return sw, true
	}
	return nil, false
}
