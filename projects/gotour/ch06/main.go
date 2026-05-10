package main

import (
	"fmt"
	"io"
)

// 定义一个结构体
type Person struct {
	Name string
	Age  int
	Addr address

	// 新增：内部数据缓冲区，作为读写的数据源
	buffer      []byte    // 存储数据的缓冲区
	readPos     int       // 当前读位置（模拟游标）
}

type address struct {
	Province string
    City string
}

// 定义一个接口
type Stringer interface {
	String() string //调用者可以通过它的 String() 方法获取一个字符串
}

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}
//Go 语言中没有继承的概念，所以结构、接口之间也没有父子关系，Go 语言提倡的是组合，利用组合达到代码复用的目的
type ReadWriter interface {
	Reader
	Writer
}

/**
 * 结构体与接口
 */
func main() {
	var p Person
	p.Name = "Alice"
	p.Age = 30
	//p := Person{Name: "Alice", Age: 30} | p:=Person{"Alice",30}

	p.Addr.Province = "Beijing"
	p.Addr.City = "Beijing"//p.Addr = address{Province: "Beijing", City: "Beijing"}
	fmt.Println(p)

	fmt.Println(p.String())
	
	//以值类型接收者实现接口的时候，不管是类型本身，还是该类型的指针类型，都实现了该接口
	printStringer(p)
	printStringer(&p)
	
	//以指针类型接收者实现接口的时候，只有该类型的指针类型实现了该接口,以值类型接收者实现接口会报错
	//printStringer(p.Addr) "报错"
	printStringer(&p.Addr)

	//工厂函数，自定义结构体，返回一个指针类型
	p1 := NewPerson("Bob", 20)
	printStringer(p1)

	//工厂函数，返回一个error接口，其实具体实现是*errorString
	errorString := New("error")
	fmt.Println(errorString)

	//创建切片
	p1.buffer = []byte("HelloWorld")
	p1.readPos = 0

	buff := make([]byte, 1024)        // 创建接收缓冲区

	n,err := p1.Read(buff)           // 从 p1 的 buffer 读入 buff
	fmt.Println(n, err)
	n,err = p1.Write(buff)
	fmt.Println(n,err)
	n,err = p1.ReadWrite(buff)
	fmt.Println(n,err)

	//接口变量 s 称为接口 fmt.Stringer 的值，它被 p1 赋值
	var s fmt.Stringer
	s = p1
	
	//类型断言表达式 s.(*Person)，
	p2, ok := s.(*Person)
	if ok {
		fmt.Println(p2)
	} else {
		fmt.Println("s is not a *Person")
	}

	//这个代码在编译的时候不会有问题，因为 address 实现了接口 Stringer，但是在运行的时候，会抛出如下异常信息：
	// panic: interface conversion: fmt.Stringer is *main.person, not main.address
	a,ok := s.(*address)
	if ok {
        fmt.Println(a)
    }else {
        fmt.Println("s不是一个address")
    }
}

//给结构体类型 person 定义一个方法,结构体 person 就实现了 Stringer 接口
func (p Person) String() string {
	return fmt.Sprintf("%s is %d years old", p.Name, p.Age)
}

//给结构体类型 address 定义一个方法,结构体 address 就实现了 Stringer 接口
func (addr *address) String() string {
	return fmt.Sprintf("the addr is %s %s", addr.Province, addr.City)
}

//定义一个可以打印 Stringer 接口的函数
func printStringer(s Stringer) {
	fmt.Println(s.String())
}

//工厂函数一般用于创建自定义的结构体
func NewPerson(name string, age int) *Person {
	return &Person{Name: name, Age: age}
}

//工厂函数，返回一个error接口，其实具体实现是*errorString
func New(text string) error {
	return &errorString{text}
}

//结构体，内部一个字段s，存储错误信息
type errorString struct {
	s string
}

//用于实现error接口
func (e *errorString) Error() string {
	return e.s
}

func (p *Person) Read(buf []byte) (n int, err error) {
	// 检查是否还有数据可读
	if p.readPos >= len(p.buffer) {
		return 0, io.EOF   // 数据读完，返回 EOF
	}
	
	// 计算本次能读多少字节
	// min(剩余数据大小, 缓冲区容量)
	remaining := len(p.buffer) - p.readPos
	toRead := len(buf)

	fmt.Printf("remaining: %d bytes, toRead:	%d bytes\n", remaining, toRead)

	if remaining < toRead {
		toRead = remaining
	}

	
	// 将数据从内部 buffer 复制到传入的 buf
	copy(buf, p.buffer[p.readPos : p.readPos+toRead])
	
	fmt.Printf("buf: %s\n", string(buf[:toRead]))

	// 推进读位置
	p.readPos += toRead
	
	return toRead, nil   // 返回实际读取的字节数
}

func (p *Person) Write(data []byte) (n int, err error) {
	// 将传入的数据追加到内部 buffer
	p.buffer = append(p.buffer, data...)
	
	return len(data), nil   // 返回实际写入的字节数
}

func (p *Person) ReadWrite(buf []byte) (n int, err error) {
	// 先读取数据（从内部 buffer 读到 buf）
    readN, readErr := p.Read(buf)
    
    // 再写回数据（从 buf 追加到内部 buffer）
    writeN, writeErr := p.Write(buf[:readN])
    
    // 合并返回值：总字节数 + 错误处理
    n = writeN
    if readErr != nil {
        err = readErr          // 优先返回读错误（如 io.EOF）
    } else {
        err = writeErr         // 否则返回写错误
    }
    
    return n, err
}