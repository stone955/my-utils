package timewheel

import (
	"container/list"
	"errors"
	"time"
)

type TimeWheel interface {
	Start()
	Stop()
	AddTimer(time.Duration, interface{}, interface{}) error
	RemoveTimer(key interface{}) error
}

// Job 延迟任务回调函数
type CallBack func(interface{})

// TimeWheel 时间轮
type timeWheel struct {
	interval time.Duration // 指针每个多长时间移动一格
	ticker   *time.Ticker
	slots    []*list.List // 时间轮的槽，每个槽保存一个双向链表，用来存放延迟任务
	// key: 定时器唯一标识 value: 定时器所在的槽, 主要用于删除定时器, 不会出现并发读写，不加锁直接访问
	timer             map[interface{}]int
	currPos           int              // 当前指针指向的槽位
	slotNum           int              // 槽的数量
	callback          CallBack         // 定时器回调函数
	addTaskChannel    chan *task       // 新增任务channel
	removeTaskChannel chan interface{} // 删除任务channel
	stopChannel       chan struct{}    // 停止定时器channel
}

// Task 延迟任务
type task struct {
	delay  time.Duration // 延迟时间
	circle int           // 时间轮需要转动几圈
	key    interface{}   // 定时器唯一标识，用于删除定时器
	data   interface{}   // 回调函数参数
}

func New(interval time.Duration, slotNum int, callback CallBack) (TimeWheel, error) {
	if interval <= 0 || slotNum <= 0 || callback == nil {
		return nil, errors.New("interval, slotNum, callback invalid")
	}

	tw := timeWheel{
		interval:          interval,
		slots:             make([]*list.List, slotNum),
		timer:             make(map[interface{}]int),
		currPos:           0,
		slotNum:           slotNum,
		callback:          callback,
		addTaskChannel:    make(chan *task),
		removeTaskChannel: make(chan interface{}),
		stopChannel:       make(chan struct{}),
	}

	tw.initSlots()

	return &tw, nil
}

// 初始化槽位，为每个槽位创建一个双向链表
func (tw *timeWheel) initSlots() {
	for i := 0; i < len(tw.slots); i++ {
		tw.slots[i] = list.New()
	}
}

// Start 启动时间轮，初始化ticker
func (tw *timeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.run()
}

// Stop 停止时间轮
func (tw *timeWheel) Stop() {
	tw.stopChannel <- struct{}{}
}

func (tw *timeWheel) AddTimer(delay time.Duration, key interface{}, data interface{}) error {
	if delay < 0 {
		return errors.New("delay can't less than zero")
	}
	if key == nil {
		return errors.New("key can't be nil")
	}
	task := task{
		delay: delay,
		key:   key,
		data:  data,
	}
	tw.addTaskChannel <- &task
	return nil
}

func (tw *timeWheel) RemoveTimer(key interface{}) error {
	if key == nil {
		return errors.New("key can't be nil")
	}
	tw.removeTaskChannel <- key
	return nil
}

// 监听各个 channel 处理接收到的消息
func (tw *timeWheel) run() {
	for {
		select {
		case <-tw.ticker.C:
			tw.handleTicker()
		case task := <-tw.addTaskChannel:
			tw.addTask(task)
		case key := <-tw.removeTaskChannel:
			tw.removeTask(key)
		case <-tw.stopChannel:
			tw.ticker.Stop()
			return
		}
	}
}

// 指针在当前槽位的处理逻辑
func (tw *timeWheel) handleTicker() {
	//处理当前槽位的延迟任务
	l := tw.slots[tw.currPos]
	tw.scanAndRunTask(l)
	// 移动到下一个槽位
	tw.moveTicker()
}

// 扫描槽位中的过期定时器，并执行回调函数
func (tw *timeWheel) scanAndRunTask(l *list.List) {
	for e := l.Front(); e != nil; {
		task := e.Value.(*task)
		// 处理本轮不执行的延迟任务
		if task.circle > 0 {
			task.circle--
			e = e.Next() // 指向下一个任务
			continue
		}
		// 处理本轮执行的延迟任务
		go tw.callback(task.data)
		next := e.Next() // 指向下一个任务
		l.Remove(e)      // 移除执行的任务
		if task.key != nil {
			delete(tw.timer, task.key)
		}
		e = next
	}
}

// 移动指针指向的槽位
func (tw *timeWheel) moveTicker() {
	if tw.currPos == tw.slotNum-1 {
		tw.currPos = 0
	} else {
		tw.currPos++
	}
}

// 添加延迟任务
func (tw *timeWheel) addTask(task *task) {
	// 计算延迟任务应该存放的槽位和时间轮需要转的圈数
	pos, circle := tw.getPosAndCircle(task.delay)
	task.circle = circle
	tw.slots[pos].PushBack(task)
	tw.timer[task.key] = pos
}

func (tw *timeWheel) getPosAndCircle(delay time.Duration) (pos, circle int) {
	// 统一时间轮单位和延迟任务时间单位
	delaySeconds := int(delay.Seconds())
	intervalSeconds := int(delay.Seconds())
	circle = delaySeconds / intervalSeconds / tw.slotNum
	pos = tw.currPos + delaySeconds/intervalSeconds%tw.slotNum
	return
}

// 移除延迟任务
func (tw *timeWheel) removeTask(key interface{}) {
	// 获取定时器所在的槽
	pos, ok := tw.timer[key]
	if !ok {
		return
	}
	// 获取槽指向的延迟任务链表
	l := tw.slots[pos]
	for e := l.Front(); e != nil; {
		task := e.Value.(*task)
		if key == task.key {
			delete(tw.timer, key)
			l.Remove(e)
		}
		e = e.Next()
	}
}
