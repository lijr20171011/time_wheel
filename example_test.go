package timewheel_test

import (
	"fmt"
	"log"
	"my/code_test/sometest/drafts/test/test20_timewheel/timewheel"
	"my/printl"
	"testing"
	"time"
)

func TestTimeWheel(t *testing.T) {
	printl.Info("====")
	tsk, err := timewheel.NewTask("println", timewheel.SingleTask, 1, 5*time.Second, Println, "aaa")
	if err != nil {
		log.Fatal("生成任务异常;", err)
	}
	printl.Info("====")
	tw := timewheel.NewTimeWheel()
	err = tw.AddTask(tsk)
	if err != nil {
		log.Fatal("添加任务至时间轮异常;", err)
	}
	printl.Info("====")
	time.Sleep(15 * time.Second)
}

func Println(msg ...interface{}) (err error) {
	fmt.Println(time.Now(), msg)
	return nil
}
