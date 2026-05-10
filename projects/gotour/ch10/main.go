package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Context 接口只有四个方法：Deadline、Done、Err、Value
type Context interface {
	//获取设置的截止时间，deadline 是截止时间，ok表示是否设置了截止时间
	Deadline() (deadline time.Time, ok bool)
	//Done（常用） 方法返回一个只读的 channel，类型为空结构体，用于通知停止
	Done() <-chan struct{} 
	//获取错误信息，如果有错误，返回错误信息，否则返回nil
	Err() error
	//获取该 Context 上绑定的值（键值对），key 是键，value 是值
	Value(key interface{}) interface{}
}

/**
 * 多线程并发控制神器 Context 
 */
func main() {
	fmt.Println("===协程如何退出===")
	fmt.Println("===select+channel 方案===")
	//使用 select+channel，通过 channel 发送指令让监控狗停止，进而达到协程退出的目的
	var wg sync.WaitGroup
	wg.Add(1)
	stopCh := make(chan bool) //用来停止监控狗
	go func() {
		defer wg.Done()
		watchDog(stopCh,"【监控狗1】")
	}()
	time.Sleep(5 * time.Second) //先让监控狗监控5秒
	stopCh <- true //发停止指令
	wg.Wait()

	fmt.Println("===Context 实现===")
	// Background()：空 Context，不可取消，没有截止时间，主要用于 Context 树的根节点。
	// WithCancel(parent Context)：可取消的 Context：用于发出取消信号，当取消的时候，它的子 Context 也会取消。
	// WithDeadline(parent Context, d time.Time)：可定时取消的 Context：多了一个定时的功能。
	// WithTimeout(parent Context, timeout time.Duration)：可超时取消的 Context：多了一个超时取消的功能。
	// WithValue(parent Context, key, val interface{})：值 Context：用于存储一个 key-value 键值对。

	fmt.Println("===Context 方案===")
	// Context控制多个协程之间的协作，尤其是取消操作
	wg.Add(3)
	// 使用Context API，通过WithCancel创建一个可以取消的上下文
	// context.Background() 用于生成一个空 Context，一般作为整个 Context 树的根节点
	ctx,stop:=context.WithCancel(context.Background())
	// 使用WithCancel创建子Context
	// ctx1,stop1:=context.WithCancel(ctx)
	// ctx2,stop2:=context.WithCancel(ctx)

	// 取消多个协程也比较简单，把 Context 作为参数传递给协程即可
	go func() {
		defer wg.Done()
		watchDog1(ctx,"【监控狗1】")
	}()
	
	go func() {
		defer wg.Done()
		watchDog1(ctx,"【监控狗2】")
	}()

	go func() {
		defer wg.Done()
		watchDog1(ctx,"【监控狗3】")
	}()
	time.Sleep(5 * time.Second) //先让监控狗监控5秒
	// context.WithCancel 函数返回的取消函数 通知goroutine退出
	// 如果一个 Context 有子 Context，在该 Context 取消时，它的子Context也会取消，父context不会取消。
	stop()
	time.Sleep(time.Second) 

	fmt.Println("===Context 传值===")

	// 创建新的context
	userCtx, cancelUser := context.WithCancel(context.Background())
	// 给context绑定值
	valCtx := context.WithValue(userCtx, "userId", 2)
	wg.Add(1)
	go func() {
		defer wg.Done()
		getUser(valCtx)
	}()
	time.Sleep(3 * time.Second) // 给 getUser 执行时间
	cancelUser()                 // 取消 userCtx
	wg.Wait()
	fmt.Println("程序结束")

	fmt.Println("===Context 使用原则===")
	// Context 不要放在结构体中，要以参数的方式传递。
	// Context 作为函数的参数时，要放在第一位，也就是第一个参数。
	// 要使用 context.Background 函数生成根节点的 Context，也就是最顶层的 Context。
	// Context 传值要传递必须的值，而且要尽可能地少，不要什么都传。
	// Context 多协程安全，可以在多个协程中放心使用。

	fmt.Println("===Context 实现日志跟踪===")
	logCtx := NewRequestContext(context.Background(), "req-12345", "user-001")
	processOrder(logCtx)
}

// 日志跟踪 Key 类型（避免冲突）
type traceKey struct{}
type userKey struct{}
type requestIDKey struct{}

// 请求上下文结构
type RequestContext struct {
	context.Context
	RequestID string
	UserID    string
}

// 创建带跟踪信息的 Context
func NewRequestContext(parent context.Context, requestID, userID string) *RequestContext {
	ctx := context.WithValue(parent, requestIDKey{}, requestID)
	ctx = context.WithValue(ctx, userKey{}, userID)
	return &RequestContext{Context: ctx, RequestID: requestID, UserID: userID}
}

// 从 Context 获取日志信息
func getTraceInfo(ctx context.Context) (requestID, userID string) {
	if v := ctx.Value(requestIDKey{}); v != nil {
		requestID = v.(string)
	}
	if v := ctx.Value(userKey{}); v != nil {
		userID = v.(string)
	}
	return
}

// 日志记录函数
func log(ctx context.Context, msg string) {
	requestID, userID := getTraceInfo(ctx)
	fmt.Printf("[TRACE] requestID=%s userID=%s | %s\n", requestID, userID, msg)
}

// 处理订单（模拟多个 goroutine 协作）
func processOrder(ctx context.Context) {
	log(ctx, "开始处理订单")
	var wg sync.WaitGroup
	wg.Add(2)

	// 验证库存
	go func() {
		defer wg.Done()
		log(ctx, "验证库存中...")
		time.Sleep(500 * time.Millisecond)
		log(ctx, "库存验证通过")
	}()

	// 扣减余额
	go func() {
		defer wg.Done()
		log(ctx, "扣减余额中...")
		time.Sleep(300 * time.Millisecond)
		log(ctx, "余额扣减成功")
	}()

	wg.Wait()
	log(ctx, "订单处理完成")
}

func watchDog(stopCh chan bool,name string) {
	//开启for select循环，一直后台监控
	for{ //for无限循环
		select {
		case <-stopCh:
			fmt.Println(name,"停止指令已收到，马上停止")
			return
		default:
			fmt.Println(name,"正在监控……")
		}
		time.Sleep(1*time.Second)
	}
}

func watchDog1(ctx context.Context,name string) {
	//开启for select循环，一直后台监控
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name,"停止指令已收到，马上停止")
			return
		default:
			fmt.Println(name,"正在监控……")
		}
		time.Sleep(1 * time.Second)
	}
}

func getUser(ctx context.Context){
	for  {
		select {
		case <-ctx.Done():
			fmt.Println("【获取用户】","协程退出")
			return
		default:
			userId:=ctx.Value("userId")
			fmt.Println("【获取用户】","用户ID为：",userId)
			time.Sleep(1 * time.Second)
		}
	}
}