package main

import (
	"fmt"
	"sync"
	"time"
)

var(
   sum int
   mutex sync.Mutex
   rwMutex sync.RWMutex
)

/**
 * 并发编程
 * 1. 资源竞争问题
 * 2. 互斥锁 Mutex
 * 3. 读写锁 RWMutex
 * 4. 协程完成 WaitGroup
 * 5. 单例 Once
 * 6. 发号施令 Cond
 * 7. 线程安全 Map
 */
func main() {
	// 资源竞争问题
	// 开启100个协程让sum+10
	// 多个协程交叉执行 sum+=i，产生不可预料的结果（结果可能不是1000）。
	for i := 0; i < 100; i++ {
		go add(10)
	}
	//防止提前退出
	time.Sleep(2 * time.Second)
	fmt.Println("和为:",sum)

	fmt.Println("===互斥锁 Mutex===")
	// 加锁解决资源竞争问题
	sum = 0
	for i := 0; i < 100; i++ {
		go add1(10)
	}
	time.Sleep(2 * time.Second)
	fmt.Println("和为:",sum)

	/*
	//add 和 readSum 使用同一个 sync.Mutex（性能差）
	for i := 0; i < 100; i++ {
		go add(10)
	}
	for i:=0; i<10;i++ {
		go fmt.Println("和为:",readSum())
	}
	time.Sleep(2 * time.Second)
	*/

	fmt.Println("===读写锁 RWMutex===")
	// 多个 goroutine 可以同时读数据，不再相互等待
	sum = 0
	for i := 0; i < 100; i++ {
		go add1(10)
	}
	for i:=0; i<10;i++ {
		go fmt.Println("和为:",readSum1())
	}
	time.Sleep(2 * time.Second)

	fmt.Println("===协程完成 WaitGroup===")
	// 使用 sync.WaitGroup 等待所有 goroutine 完成
	run()

	fmt.Println("===单例 Once===")
	// 使用 sync.Once 确保只执行一次
	doOnce()

	fmt.Println("===发号施令 Cond===")
	// 使用 sync.Cond 实现发号施令
	race()

	fmt.Println("===线程安全 Map===")
	// 使用 sync.Map 实现线程安全的 map
	// sync.Map 适用于读多写少场景
	// 普通 map 并发读写需加 sync.Mutex 或 sync.RWMutex
	// 不要混用 sync.Map 和普通 map
	sMap()
}

func add(i int) {
   	sum += i
}

func add1(i int) {
   	mutex.Lock() // 加锁
	defer mutex.Unlock() // 解锁（采用 defer，确保锁一定会被释放）

   	// 被加锁保护的 sum+=i 代码片段又称为临界区
	// 是一个访问共享资源的程序片段，当有协程进入临界区段时，其他协程必须等待
	sum += i
}

// 增加了一个读取sum的函数，便于演示并发
func readSum() int {
	mutex.Lock()
    defer mutex.Unlock()
	b:=sum
	return b
}

func readSum1() int {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	b:=sum
	return b
}

func run () {
	sum = 0

	// 创建一个WaitGroup用于等待所有goroutine完成 
	var wg sync.WaitGroup
	
	// 设置计数器的值
	// 因为要监控110个协程，所以设置计数器为110
	wg.Add(110)
	for i := 0; i < 100; i++ {
      	go func() {
			//计数器值减1
			defer wg.Done()
			add1(10)
      	}()
	}
	for i:=0; i<10;i++ {
		go func() {
			//计数器值减1
			defer wg.Done()
			fmt.Println("和为:",readSum1())
		}()
	}
	//一直等待，只要计数器值为0
	wg.Wait()
}

func doOnce() {
	// 创建一个sync.Once实例
	var once sync.Once
	onceBody := func() {
		fmt.Println("Only once")
	}
	//用于等待协程执行完毕
	done := make(chan bool)
	//启动10个协程执行once.Do(onceBody)
	for i := 0; i < 10; i++ {
		go func() {
			//把要执行的函数(方法)作为参数传给once.Do方法即可
			once.Do(onceBody)
			done <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}

//10个人赛跑，1个裁判发号施令
func race(){
	//创建一个Cond实例，并传入一个锁对象
	cond :=sync.NewCond(&sync.Mutex{})
	var wg sync.WaitGroup
	wg.Add(11)
	for i:=0;i<10; i++ {
		go func(num int) {
			defer  wg.Done()
			fmt.Println(num,"号已经就位")
			cond.L.Lock() //加锁
			cond.Wait()//等待发令枪响 // 当前协程进入等待状态，并释放锁（必须在加锁后才能调用）
			fmt.Println(num,"号开始跑……")
			cond.L.Unlock() //解锁
		}(i)
	}
	//等待所有goroutine都进入wait状态
	time.Sleep(2*time.Second)
	go func() {
		defer  wg.Done()
		fmt.Println("裁判已经就位，准备发令枪")
		fmt.Println("比赛开始，大家准备跑")
		cond.Broadcast()//发令枪响 // 所有等待的协程都会被唤醒
		//cond.Signal() // 唤醒一个等待时间最长的协程
		// 注意：在调用 Signal 或者 Broadcast 之前，要确保目标协程处于 Wait 阻塞状态，不然会出现死锁问题
	}()
	//防止函数提前返回退出
	wg.Wait()
}

// Map 的操作方法
// Store：存储一对 key-value 值。
// Load：根据 key 获取对应的 value 值，并且可以判断 key 是否存在。
// LoadOrStore：如果 key 对应的 value 存在，则返回该 value；如果不存在，存储相应的 value。
// Swap：替换 key 对应的 value，并返回原始 value。
// Delete：删除一个 key-value 键值对。
// Range：循环迭代 sync.Map，效果与 for range 一样
func sMap() {
	// 创建一个 sync.Map 实例
	var m sync.Map

	// 存储键值对
	m.Store("name", "飞雪无情")
	m.Store("age", 30)

	// 读取值
	value, ok := m.Load("name")
	fmt.Println(value, ok) // 飞雪无情 true

	// 读取或删除
	m.Delete("name")

	// 遍历
	m.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true // 返回 false 停止遍历
	})

	// LoadOrStore: 读取已存在的值，如果不存在则写入
	v, loaded := m.LoadOrStore("name", "新名字")
	fmt.Println(v, loaded) //loaded：true = key 已存在（读取）；false = key 不存在（写入新值）

	// Swap: 替换并返回旧值（强制覆盖）
	old, _ := m.Swap("name", "替换后的名字")
	fmt.Println(old) 

	v, _ = m.Load("name")
	fmt.Println(v)   // "新名字"（已成功替换）
}