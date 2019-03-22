# time_wheel
尝试时间轮定时器

# 使用

## 新建时间轮

```go
    tw := timewheel.NewTimeWheel()
```
## 新建任务
```go
    tsk, err := timewheel.NewTask(taskName, taskType, execCnt, interval, callFunc, callArgs)
```

### 参数说明
* `taskName` : 任务名称，在同一`timeWheel`中唯一
* `taskType` : 任务类型，可以取以下值
    * `SingleTask` : 单次任务
    * `MultiTask` : 多次任务
    * `PermanentTask` : 永久任务
* `execCnt` : 任务执行次数，对单次任务和永久任务没有影响
* `interval` : 任务执行时间间隔
* `callFunc` : 执行的任务 
    * 任务格式 : `func(args ...interface{}) error`
* `callArgs` : 任务参数

## 添加任务至时间轮
```go
    err = tw.AddTask(tsk)
```

## 启动时间轮
```go
    tw.Start()
```

## 停止时间轮
```go
    tw.Stop()
```

## 判断时间轮中任务是否存在
```go
	ok := tw.IsTaskExist(taskName) // taskName : 任务名称
```

# 测试

## 代码

```go
func main() {
	taskName := "println"
	// 新建时间轮
	tw := timewheel.NewTimeWheel()
	fmt.Println("new time wheel")
	// 新建任务
	tsk, err := timewheel.NewTask(taskName, timewheel.MultiTask, 2, 5*time.Second, Println, "aaa")
	if err != nil {
		log.Fatal("生成任务异常;", err)
	}
	fmt.Println("new task")
	// 添加任务至时间轮
	err = tw.AddTask(tsk)
	if err != nil {
		log.Fatal("添加任务至时间轮异常;", err)
	}
	fmt.Println("add task")
	// 启动时间轮
	tw.Start()
	fmt.Println("time wheel started")

	// ==== 测试
	fmt.Println(tw.IsTaskExist(taskName)) // true

	time.Sleep(6 * time.Second)

	fmt.Println(tw.IsTaskExist(taskName)) // true

	time.Sleep(5 * time.Second)

	fmt.Println(tw.IsTaskExist(taskName)) // false

	tw.Stop()
	fmt.Println("time wheel stoped")
}

func Println(msg ...interface{}) (err error) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), msg)
	return nil
}

```

## 预计输出

```shell
new time wheel
new task
add task
time wheel started
true
2019-03-22 17:38:09 [aaa]
true
2019-03-22 17:38:14 [aaa]
false
time wheel stoped
```
