package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
	"unsafe"
)

// 与 runtime 中切片头布局一致（教学用途）
type sliceHeader struct {
	Data uintptr // 底层数组首地址
	Len  int     // 当前长度
	Cap  int     // 当前容量
}

func main() {
	fmt.Println("===切片共享底层数组===")
	a1 := [2]string{"飞雪无情", "张三"}
	s1 := a1[0:1]
	s2 := a1[:]
	//打印出s1和s2的Data值，是一样的，说明切片s1和s2共享底层数组
	//SliceHeader是切片在程序运行时的真实结构
	fmt.Println((*reflect.SliceHeader)(unsafe.Pointer(&s1)).Data)
	fmt.Println((*reflect.SliceHeader)(unsafe.Pointer(&s2)).Data)

	fmt.Println("===string 和 []byte 互转===")
	//强制转换先分配一个内存再复制内容的方式，s3和s的内存地址不一样
	s := "飞雪无情"
	//stringHeader字符串在程序运行时的真实结构     取 s 的底层数据地址
	fmt.Printf("s的内存地址：%d\n", (*reflect.StringHeader)(unsafe.Pointer(&s)).Data)
	b := []byte(s)
	fmt.Printf("b的内存地址：%d\n", (*reflect.SliceHeader)(unsafe.Pointer(&b)).Data)
	s3 := string(b)
	fmt.Printf("s3的内存地址：%d\n", (*reflect.StringHeader)(unsafe.Pointer(&s3)).Data)

	b1 := []byte(s)
	// 没有申请新内存（零拷贝） 把 []byte 头直接解释成 string，而不是重新分配内存
	s4 := *(*string)(unsafe.Pointer(&b1))
	fmt.Println("s4的值为:", s4)

	//StringHeader 通过 unsafe.Pointer转为 SliceHeader ，缺少cap补上
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&s)) //SliceHeader结构体，指向s的底层数组
	sh.Cap = sh.Len                                  //将sh的容量设置为sh的长度
	b2 := *(*[]byte)(unsafe.Pointer(sh))             //将sh的底层数组转换为[]byte
	fmt.Println("b2的值为:", b2)                        //b2的值为: [102 105 120 117 101 113 110]

	fmt.Println("===sliceHeader 与切片高效处理===")

	// 1) 看清切片“头部信息”与 append 扩容行为
	demoSliceHeader()
	// 2) 对比预分配与非预分配的追加成本
	demoCapacityGrowth()
	// 3) 演示原地过滤（零额外分配）
	demoInplaceFilter()
	// 4) 演示删除区间的高效写法
	demoDeleteRange()
	// 5) 演示子切片导致的大数组被持有问题
	demoSubSliceLeak()
}

func demoSliceHeader() {
	fmt.Println("\n--- 1) sliceHeader：切片本质是三元组(Data, Len, Cap) ---")

	// 注意：切片变量本身只是一个“描述符”，并不直接存储所有元素
	s := []int{10, 20, 30, 40}
	h := (*sliceHeader)(unsafe.Pointer(&s))

	fmt.Printf("slice=%v\n", s)
	fmt.Printf("header.Data=0x%x Len=%d Cap=%d\n", h.Data, h.Len, h.Cap)
	fmt.Printf("&s[0]=%p（通常和 Data 对应）\n", &s[0])

	// append 可能触发扩容：
	// - 若容量足够：Data 通常不变，只改 Len
	// - 若容量不足：分配新数组，拷贝旧数据，Data 变化
	before := h.Data
	s = append(s, 50, 60, 70, 80, 90)
	after := (*sliceHeader)(unsafe.Pointer(&s)).Data
	fmt.Printf("append 后 Len=%d Cap=%d Data changed=%v\n", len(s), cap(s), before != after)
}

func demoCapacityGrowth() {
	fmt.Println("\n--- 2) 预分配容量：减少扩容与拷贝 ---")

	const n = 200000

	// 基准只做教学演示：比较趋势而非绝对数值（受机器和负载影响）
	start := time.Now()
	noPrealloc := []int{}
	for i := 0; i < n; i++ {
		noPrealloc = append(noPrealloc, i)
	}
	cost1 := time.Since(start)

	// 预估元素个数时，优先预分配 cap，能显著减少扩容拷贝
	start = time.Now()
	prealloc := make([]int, 0, n)
	for i := 0; i < n; i++ {
		prealloc = append(prealloc, i)
	}
	cost2 := time.Since(start)

	fmt.Printf("不预分配: len=%d cap=%d cost=%v\n", len(noPrealloc), cap(noPrealloc), cost1)
	fmt.Printf("预分配容量: len=%d cap=%d cost=%v\n", len(prealloc), cap(prealloc), cost2)
}

func demoInplaceFilter() {
	fmt.Println("\n--- 3) 原地过滤：res := src[:0] 复用底层数组 ---")

	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	// 经典原地过滤模板：
	// res 与 src 共用同一底层数组，避免重新分配
	// 适用于“保留一部分元素”的场景
	res := src[:0] // 长度归零，容量不变；可复用原底层数组
	for _, v := range src {
		if v%2 == 0 {
			res = append(res, v)
		}
	}

	fmt.Printf("过滤偶数结果=%v\n", res)
	fmt.Printf("src 与 res 共享底层数组: %t\n", &src[0] == &res[0])
}

func demoDeleteRange() {
	fmt.Println("\n--- 4) 高效删除区间：copy + 截断 ---")

	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	i, j := 3, 7 // 删除 [3,7) => 3,4,5,6
	// 时间复杂度 O(n-j)，只移动“后半段”元素
	oldLen := len(s)
	removed := j - i
	copy(s[i:], s[j:]) //将s[j:]复制到s[i:]，覆盖掉s[i:j] 拷贝后底层数组前部变成：[0 1 2 7 8 9 ...]
	// 对于包含指针元素的切片，建议 clear 掉“被丢弃的尾部”，帮助 GC
	// 这里是 int 切片，clear 主要用于演示正确范围
	clear(s[oldLen-removed : oldLen]) //清除s[oldLen-removed:oldLen]
	s = s[:oldLen-removed]

	fmt.Printf("删除区间[%d,%d)后=%v\n", i, j, s)
}

func demoSubSliceLeak() {
	fmt.Println("\n--- 5) 子切片内存保留陷阱：小切片引用大数组 ---")

	big := make([]byte, 8<<20) // 8MB
	for i := range big {
		big[i] = byte(i)
	}
	smallView := big[:16] // 只要 smallView 活着，8MB 大数组通常也不能被回收
	fmt.Printf("smallView len=%d cap=%d\n", len(smallView), cap(smallView))

	// 正确做法：
	// 将小片段 copy 到新切片，切断对 big 底层数组的引用
	// 常用于“从大响应体里截取小字段”这类场景
	smallCopy := append([]byte(nil), smallView...)
	big = nil
	runtime.GC()

	fmt.Printf("smallCopy len=%d cap=%d（独立内存）\n", len(smallCopy), cap(smallCopy))
}
