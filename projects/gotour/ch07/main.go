package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

//自定义 error
type CommonError struct {
	Code int
	Message string
}

//错误嵌套
type MyError struct {
	err error
	msg string
}

/**
 * 错误处理
 * error 处理和断言
 * defer 函数
 * panic 异常
 * recover 函数
 */
func main() {
	i,err := strconv.Atoi("a")//字符串转整数
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i)
	}

	sum,err := add(-1,2)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(sum)
	}

	fmt.Println("===error 断言===")
	//类型断言检查是否为CommonError类型
	if cm,ok:=err.(*CommonError);ok{
		fmt.Println("错误代码为:",cm.Code,"，错误信息为：",cm.Message)
	} else {
		fmt.Println(sum)
	}

	fmt.Println("===错误嵌套===")
	newErr := MyError{err, "数据上传问题"}
	fmt.Println(newErr.Error())

	fmt.Println("===Error Wrapping 功能===")
	e := errors.New("原始错误e")
	w := fmt.Errorf("Wrap了一个错误:%w", e)
	fmt.Println(w)

	fmt.Println("===errors.Unwrap 函数===")
	//用于获取被嵌套的 error
	fmt.Println(errors.Unwrap(w))

	fmt.Println("===errors.Is 函数===")
	//用于检查 error 是否为某个特定的 error
	fmt.Println(errors.Is(w,e))

	fmt.Println("===errors.As 函数===")
	//有了 error 嵌套后，error 断言也不能用了,前面 error 断言的例子，可以使用 errors.As 函数重写
	var cm *CommonError
	if errors.As(err,&cm){
		fmt.Println("错误代码为:",cm.Code,"，错误信息为：",cm.Message)
	} else {
		fmt.Println(sum)
	}

	fmt.Println("===Deferred 函数===")
	// 注册时机	defer 执行时只是注册，不是立即执行
	// 执行时机	函数返回前按 LIFO（后进先出）顺序执行
	// 位置要求	必须在 return 之前才会被执行到
	// 参数求值	注册时立即求值，不是执行时求值
	read,err := ReadFile("./main.go")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(read))
	}

	fmt.Println("=== 多个defer开始 ===")
	fmt.Println("=== 倒序执行，先进后出: ===")
	//引用的方法中的defer相当于一个整体与当前方法中的defer按照倒序执行，先进后出
    test()
	fmt.Println("=== 参数求值陷阱: ===")
	demo()
	//demoFixed()
    fmt.Println("=== 结束 ===")

	fmt.Println("===Panic 异常===")
	//运行时的问题会引起 panic 异常,也可以手动抛出 panic 异常
	//panic 异常是一种非常严重的情况，会让程序中断运行，使程序崩溃，所以如果是不影响程序运行的错误，不要使用 panic，使用普通错误 error 即可。
	//connectMySQL("","root","root")//致命错误，程序崩溃！
	
	fmt.Println("===Recover 函数：捕获 Panic 异常===")
	//recover 函数用于恢复 panic 异常，只有在 defer 函数中调用才有用
	//defer 关键字 + 匿名函数 + recover 函数从 panic 异常中恢复的方式
	//recover 函数返回的值就是通过 panic 函数传递的参数值
	defer func() {//defer延迟执行，在函数返回前触发, func() {} 匿名函数（闭包）
		if err := recover(); err != nil { //recover()内置函数，捕获 panic传递的值
			fmt.Println(err)
		}
	}() //立即执行匿名函数，返回的函数体被 defer 注册
	connectMySQL("","","")//程序恢复，不会崩溃！
}

func add(a,b int) (int,error) {
	if a < 0 || b < 0 {
		//使用 errors.New 这个工厂函数生成的错误信息
		//return 0,errors.New("a and b must be non-negative")

		return 0,&CommonError{Code: 1001,Message: "a and b must be non-negative"}
	}
	return a + b,nil
}

//自定义 error
func (e *CommonError) Error() string {
	return fmt.Sprintf("Error code: %d, Error message: %s", e.Code, e.Message)
}

func (e *MyError) Error() string {
	 return e.err.Error() + e.msg
}

//自定义读取文件方法（ioutil：1.16+ 已废弃，所有函数已迁移到 io 和 os 包中）
func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()  // ✅ 确保函数退出时关闭文件

    return io.ReadAll(f)
	//关键点：
	// os.Open() 返回的 *File 有 Close() 方法
	// defer 紧跟在打开成功之后
	// 无论函数如何退出（正常/panic），都会执行 Close()
}

//HTTP 响应体
func fetch(url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()  // ✅ 必须关闭响应体

    return io.ReadAll(resp.Body)
}

func connectMySQL(ip,username,password string) {
   if ip =="" {
      panic("ip不能为空")
   }
   //省略其他代码
}

func test() {
	defer fmt.Println("① 第一个 defer")  // 第1个注册，最后执行
    defer fmt.Println("② 第二个 defer")
    defer fmt.Println("③ 第三个 defer")
    defer fmt.Println("④ 第四个 defer")  // 最后注册，最先执行
	fmt.Println("函数自身代码")
}

func demo() {
    for i := 0; i < 3; i++ {
        defer func() {
            fmt.Println(i)  // 引用外部变量 i
        }()  // 注意：i 是引用！
    }

    // 等价于打印三次 i 的最终值
    // 输出: 2 2 2 (循环结束后 i=2)
	// 在 Go 1.22+ 版本中，for 循环的每次迭代都会创建新的变量实例，
	// 所以三个 defer 捕获的是不同的 i（值分别为 0, 1, 2），再按 LIFO 顺序输出就是 2 1 0
	// < 1.22 版本中，for 循环的每次迭代都会复用同一个变量实例，所以三个 defer 捕获的是同一个变量 i，最终值为 2
}

func demoFixed() {
    for i := 0; i < 3; i++ {
        defer func(n int) {  // ✅ 通过参数传入
            fmt.Println(n)
        }(i)  // 注册时就求值！
    }

    // 输出: 2 1 0 (每次传入当时的值)
}
