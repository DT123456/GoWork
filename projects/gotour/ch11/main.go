package main

import (
	"fmt"
	"sync"
	"time"
)

/**
 * 高效并发模式
 * for select循环、
 * for range select 有限循环模式
 * select timeout 超时模式
 * Pipeline流水线
 * 扇出扇入模式
 * Futures 模式
 * Semaphore 信号量模式
 */
func main() {
	fmt.Println("===for select 循环模式===")
	//for select 无限循环模式 参考ch10/main.go watchDog\ 

	fmt.Println("===for range select 有限循环模式===")
	done:=make(chan bool)
	resultCh:=make(chan int)

	go func() {
		for i := 1; i <= 3; i++ {
		time.Sleep(time.Second)
		resultCh <- i  // 发送数据
		}
		done <- true
	}()

	loop : for _,v:=range []int{1,2,3,4}{
		select {
		case <- done:
			fmt.Println("完成", "第", v, "次")
			break loop // 退出 label 标记的 for 循环，单break不加loop表示退出select，使用return会退出main，后续不执行
		case s := <- resultCh:
			fmt.Println("收到结果:", s, "第", v, "次")
			time.Sleep(time.Second)
		}
	}

	fmt.Println("准备退出循环")
	fmt.Println("===select timeout 模式===")
	result := make(chan string)
	go func() {
		//模拟网络访问
		time.Sleep(8 * time.Second)
		result <- "服务端结果"
	}()
	select {
	case v := <-result:
		fmt.Println(v)
	case <-time.After(5 * time.Second): // 优先使用 Context 的 WithTimeout 函数超时取消
		fmt.Println("网络访问超时了")
	}

	fmt.Println("===Pipeline 模式===")
	//流水线模式
	coms := buy(10)    //采购10套配件
	phones := build(coms) //组装10部手机
	packs := pack(phones) //打包它们以便售卖
	//输出测试，看看效果
	for p := range packs {
		fmt.Println(p)
	}

	fmt.Println("===扇出和扇入模式===")
	coms1 := buy(100)    //采购100套配件
	//三班人同时组装100部手机
	phones1 := build(coms1)
	phones2 := build(coms1)
	phones3 := build(coms1)
	//汇聚三个channel成一个
	phonesMerge := merge(phones1,phones2,phones3)
	packs1 := pack(phonesMerge) //打包它们以便售卖
	//输出测试，看看效果
	for p := range packs1 {
		fmt.Println(p)
	}

	fmt.Println("===Futures 模式===")
	//未来模式 多个任务并行执行，最后汇总结果
	vegetablesCh := washVegetables() //洗菜
	waterCh := boilWater()           //烧水
	fmt.Println("已经安排洗菜和烧水了，我先眯一会")
	time.Sleep(2 * time.Second)

	fmt.Println("要做火锅了，看看菜和水好了吗")
	vegetables := <-vegetablesCh
	water := <-waterCh
	fmt.Println("准备好了，可以做火锅了:",vegetables,water)

	fmt.Println("===Semaphore（信号量）===")
	// 控制并发数量
	// 典型应用场景
	// 限流：控制 API 调用频率
	// 连接池：限制数据库/Redis 连接数
	// 资源池：限制线程数、文件句柄数
	fmt.Println("===模拟并发请求控制：最多同时3个请求===")
	limit := 3
	sem := make(chan struct{}, limit) //创建带缓冲 channel，容量=3
	urls := []string{"url1", "url2", "url3", "url4", "url5", "url6", "url7", "url8", "url9", "url10"}

	var wg sync.WaitGroup
	for _, url := range urls {
		sem <- struct{}{} // 获取信号量（阻塞直到有空闲位置）
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			fmt.Printf("[开始] %s\n", u)
			time.Sleep(time.Second)
			fmt.Printf("[完成] %s\n", u)
			<-sem // 释放信号量
		}(url)
	}
	wg.Wait()
	fmt.Println("所有请求完成")
}

//工序1采购
func buy(n int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- fmt.Sprint("配件", i)
		}
	}()
	return out
}

//工序2组装
func build(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for c := range in {
			out <- "组装(" + c + ")"
		}
	}()
	return out
}

//工序3打包
func pack(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for c := range in {
			out <- "打包(" + c + ")"
		}
	}()
	return out
}

//扇入函数（组件），把多个chanel中的数据发送到一个channel中
func merge(ins ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)
	//把一个channel中的数据发送到out中
	p:=func(in <-chan string) {
		defer wg.Done()
		for c := range in {
			out <- c
		}
	}
	wg.Add(len(ins))
	//扇入，需要启动多个goroutine用于处于多个channel中的数据
	for _,cs:=range ins{
		go p(cs)
	}
	//等待所有输入的数据ins处理完，再关闭输出out
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

//洗菜
func washVegetables() <-chan string {
	vegetables := make(chan string)
    fmt.Println("开始洗菜")
	go func() {
		time.Sleep(5 * time.Second)
		vegetables <- "洗好的菜"
		fmt.Println("洗完菜了")
	}()

	return vegetables
}

//烧水
func boilWater() <-chan string {
	water := make(chan string)
	fmt.Println("开始烧水")
	go func() {
		time.Sleep(5 * time.Second)
		water <- "烧开的水"
		fmt.Println("烧好水了")
	}()

	return water
}
