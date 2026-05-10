package main

import (
	"errors"
	"fmt"
	"strings"
)

/**
 * 单元测试示例：
 * 1) 可测试的纯函数
 * 2) 错误分支处理
 * 3) 表驱动测试目标函数
 */
func main() {
	fmt.Println("===ch18: 单元测试示例===")

	score := []int{90, 80, 100}
	avg, err := Average(score)
	fmt.Println("Average:", avg, "err:", err)

	msg, err := FormatUserName("  Alice  ")
	fmt.Println("FormatUserName:", msg, "err:", err)

	fmt.Println("GradeLevel(85):", GradeLevel(85))

	fib := Fibonacci(10)
	fmt.Println("Fibonacci(10):", fib)

	// go test -v .\ch18\     // 运行所有测试用例
	// go test -v --coverprofile=ch18.cover ./ch18  // 生成覆盖率报告
	//go tool cover -html ch18.cover -o ch18.html  // 生成覆盖率报告的 HTML 文件
	// go tool cover -func ch18.cover  // 显示覆盖率报告的函数覆盖率
	// go tool cover -mode=set -func=ch18.cover  // 显示覆盖率报告的函数覆盖率 默认是 mode=set 适用于单线程环境
	// go tool cover -mode=atomic -func=ch18.cover  // 显示覆盖率报告的函数覆盖率 适用于多线程环境

	// go test -bench=. ./ch18  // 运行基准测试  -v 打印详细信息
	// go test -bench=. ./ch18 -benchtime=10s  // 运行基准测试 10秒
	// go test -bench=. ./ch18 -benchtime=10s -benchmem  // 运行基准测试 10秒 并显示内存使用情况
	// go test -bench=. ./ch18 -benchtime=10s -benchmem -cpuprofile=cpu.prof  // 运行基准测试 10秒 并显示内存使用情况 并生成 CPU 使用情况报告
	
	// go tool pprof cpu.prof  // 分析 CPU 使用情况
	// go tool pprof cpu.prof -inuse_space  // 分析 CPU 使用情况 并显示内存使用情况
	// go tool pprof cpu.prof -inuse_objects  // 分析 CPU 使用情况 并显示内存使用情况
	// go tool pprof cpu.prof -alloc_space  // 分析 CPU 使用情况 并显示内存使用情况
}

// Average 计算平均分，空切片返回错误（便于测试错误分支）。
func Average(nums []int) (float64, error) {
	if len(nums) == 0 {
		return 0, errors.New("nums cannot be empty")
	}

	sum := 0
	for _, n := range nums {
		sum += n
	}
	return float64(sum) / float64(len(nums)), nil
}

// FormatUserName 去除首尾空格，空字符串返回错误。
func FormatUserName(name string) (string, error) {
	trimmed := strings.TrimSpace(name) //去除首尾空格
	if trimmed == "" {
		return "", errors.New("name cannot be empty")
	}
	return strings.ToUpper(trimmed), nil //转换为大写
}

// GradeLevel 根据分数划分等级。
func GradeLevel(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 60:
		return "C"
	default:
		return "D"
	}
}

var cache = map[int]int{}
// Fibonacci 返回第 n 项斐波那契数列值（从 0 开始：0,1,1,2,3...）。
func Fibonacci(n int) int {
	if v, ok := cache[n]; ok {
		return v
	}
	result := 0
	switch {
	case n < 0:
		result = 0
	case n == 0:
		result = 0
	case n == 1:
		result = 1
	default:
		result = Fibonacci(n-1) + Fibonacci(n-2)
	}
	cache[n] = result
	return result
}

// FibonacciRecursive 递归实现（用于基准测试对比，不做负数校验）。
func FibonacciRecursive(n int) int {
	if n <= 1 {
		return n
	}
	return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}
