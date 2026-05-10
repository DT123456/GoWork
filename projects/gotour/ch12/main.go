package main

import "fmt"

func main() {
	fmt.Println("===指针===")
	name:="飞雪无情"
	nameP:=&name//取地址
	fmt.Println("name变量的值为:",name)
	fmt.Println("name变量的内存地址为:",nameP)

	var a int = 10
	var intP *int // 指针类型就是在对应的类型前加 * 号
	intP = &a
	fmt.Println("intP变量的值为:",intP)
	fmt.Println("intP变量的值为:",*intP)

	intP1:=new(int)  // 内置的 new 函数来声明指针类型
	fmt.Println("intP1变量的值为:",intP1)
	fmt.Println("intP1变量的内存地址为:",&intP1)
	fmt.Println("intP1变量的值为:",*intP1)

	fmt.Println("===指针的操作===")
	// 两种：一种是获取指针指向的值，一种是修改指针指向的值
	// 获取指针指向的值，使用 * 号
	nameV:=*nameP
	fmt.Println("获取nameP指针指向的值为:",nameV)
	// 修改指针指向的值，使用 * 号
	*nameP="飞雪无情2"
	fmt.Println("修改nameP指针指向的值为:",*nameP) 
	// 输出nameV：飞雪无情 nameV是独立的副本，nameP指向的值被修改了，nameV的值不变
	fmt.Println("修改nameP指针后name变量的值为:",name)

	// var 关键字直接定义的指针变量是不能进行赋值操作，没有指向的内存地址
	// 通过new 函数会申请内存空间
	var intP2 *int = new(int)
	*intP2=10
	fmt.Println("intP2变量的值为:",intP2)
	fmt.Println("intP2变量的值为:",*intP2)

	fmt.Println("===指针参数===")
	age:=18
	modifyAge(&age) // modifyAge函数使用值类型接收者，传递的是age的副本，修改的是副本的值，不会影响到age的值
	fmt.Println("age的值为:",age)

	fmt.Println("===指针接收者===")
	// 如果接收者类型是 map、slice、channel 这类引用类型，不使用指针；
	// 如果需要修改接收者，那么需要使用指针；
	// 如果接收者是比较大的类型，可以考虑使用指针，因为内存拷贝廉价，所以效率高。

	var intP3 [3]*int
	intP3[0]=&a
	intP3[1]=intP
	intP3[2]=intP1
	fmt.Println("intP3变量的值为:",intP3)
	fmt.Println("intP3变量的值为:",*intP3[0])
	fmt.Println("intP3变量的值为:",*intP3[1])
	fmt.Println("intP3变量的值为:",*intP3[2])

	fmt.Println("===引用类型===")
	m:=make(map[string]int)
	m["飞雪无情"] = 18
	fmt.Println("飞雪无情的年龄为",m["飞雪无情"])
	modifyMap(m)
	fmt.Println("飞雪无情的年龄为",m["飞雪无情"])

	//chan 同map
	//函数、接口、slice 切片都可以称为引用类型
}

func modifyAge(age *int)  {
   *age = 20
}

func modifyMap(p map[string]int)  {
    p["飞雪无情"] =20
}

