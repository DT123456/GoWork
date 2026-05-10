package main

import (
	"fmt"
	"time"
)

/**
 * 并发编程：
 * 协程（Goroutine）
 * 无缓冲Channel（通道）
 * 有缓冲Channel（通道）
 * 关闭 channel
 * 单向 channel
 * select+channel 示例
 */
func main() {
	fmt.Println("===协程（Goroutine）===")
	// Go 语言中没有线程的概念，只有协程
	// 这段代码里有两个 goroutine，一个是 main 函数启动的 main goroutine，一个是我自己通过 go 关键字启动的 goroutine

	// 启动一个新的 goroutine，go func() 启动后立即返回，不会等待函数完成
	go fmt.Println("飞雪无情")
	// 新 goroutine 与当前 goroutine 并发运行（两个 goroutine 的输出顺序不确定），何时执行新 goroutine 由 Go 调度器决定
   	fmt.Println("我是 main goroutine")
	// main 睡眠1秒，给子goroutine执行机会（没有 Sleep...，"飞雪无情"可能来不及打印就退出了！）
	// main goroutine 结束 = 整个程序退出，所有其他 goroutine 都会被强制终止！
   	time.Sleep(time.Second)
	fmt.Println("Second:", time.Second)

	fmt.Println("===Channel（通道）===")
	fmt.Println("===无缓冲Channel（通道）===")
	// chan 的操作只有两种：发送和接收
	// 数据流动、传递的场景中要优先使用 channel，它是并发安全的，性能也不错。
	ch := make(chan int)
	go func() {
		ch <- 100 // 将 100 发送到 ch 中
	}()
	fmt.Println(<- ch) // 从 ch 中接收数据

	// 程序并没有退出，可以看到"飞雪无情"的输出结果，达到了 time.Sleep 函数的效果
	// 在 main goroutine 中，从变量 ch 接收值；如果 ch 中没有值，则阻塞等待到 ch 中有值可以接收为止
	// 一个 goroutine 往管道里发送数据，另外一个从这个管道里取数据，类似于队列
	// 无缓冲 channel 的发送和接收操作是同时进行的（只传输不存储），它也可以称为同步 channel。
	ch1 := make(chan string)
	go func() {
		fmt.Println("飞雪无情")
		ch1 <- "goroutine 完成"
	}()
	fmt.Println("我是 main goroutine")
	v := <-ch1
	fmt.Println("接收到的chan中的值为：",v)

	fmt.Println("===有缓冲Channel（通道）===")
	// 有缓冲 channel 类似一个可阻塞的队列，内部的元素先进先出
	ch2 := make(chan int, 3)
	ch2 <- 1
	ch2 <- 2
	ch2 <- 3
	fmt.Println("有缓冲Channel容量为:",cap(ch2),",元素个数为：",len(ch2))
	fmt.Println(<-ch2)
	fmt.Println(<-ch2)
	fmt.Println(<-ch2)
	fmt.Println("有缓冲Channel容量为:",cap(ch2),",元素个数为：",len(ch2))

	fmt.Println("===关闭 channel===")
	close(ch2)

	fmt.Println("===单向 channel===")
	// 单向 channel 通过函数参数传递
	// counter 接收只发送的 channel，printer 接收只接收的 channel
	// Go 中无法在运行时将双向 channel 转换为单向 channel，单向 channel 只能通过函数参数传递实现。
	// 先创建双向 channel
	ch3 := make(chan int)
	// 转换为单向
	go counter(ch3) // 单向发送数据
	go printer(ch3) // 单向接收数据
	time.Sleep(time.Second) // 等待 goroutine 完成

	fmt.Println("===select+channel 示例===")
	// select 语句用于在多个 channel 上进行非阻塞的选择性接收操作
	ch4 := make(chan int)
	ch5 := make(chan int)
	//同时开启2个goroutine
	go func() {
		time.Sleep(time.Second)
		ch4 <- 100
	}()
	go func() {
		time.Sleep(time.Second * 2)
		ch5 <- 200
	}()
	// 开始select多路复用，哪个channel能获取到值，
   	// 就说明哪个最先下载好，就用哪个。
	// 同时有多个 case 可以被执行，则随机选择一个
	select {
	case v := <-ch4:
		fmt.Println("接收到ch4中的值为：",v)
	case v := <-ch5:
		fmt.Println("接收到ch5中的值为：",v)
	}
}

func counter (ch chan<- int) {
  	//函数内容使用变量out，只能进行发送操作
	for i := 0; i < 5; i++ {
		ch <- i
	}
}

func printer (ch <-chan int) {
	//函数内容使用变量in，只能进行接收操作
	for i := 0; i < 5; i++ {
		value := <-ch
		fmt.Println(value)
	}
}