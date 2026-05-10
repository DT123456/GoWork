package main

import (
	"fmt"
)

/**
 * 函数和方法
 */
func main() {
	result,err := add(1,2)
	fmt.Println(result,err)

	fmt.Println(sum(1,2,3,4,5,6,7,8,9,10))

	//匿名函数和闭包
	sum1 := func(a,b int) int {
		return a+b
	}
	fmt.Println(sum1(1,2))

	cl := closure()
	fmt.Println(cl(1))
	fmt.Println(cl(2))
	fmt.Println(cl(3))

	fmt.Println("===方法示例===")
	//方法
	age:=Age(18)
	age.String()//String() 是类型 Age 的方法
	age.Modify(20)//→ Go: (&age).Modify() Go 自动取地址: &age
	age.String()
	age.Add(1)//Go 自动取地址: &age
	age.String()

	fmt.Println("===方法赋值变量===")
	p := Person{name: "zhangsan", age: 18}
	f1 := p.Say//如果Say是方法值，赋值时就固定了 p 的值，后续的f2和f3都不会改变p的值
	f2 := p.Modify
	f3 := p.Grow
	f1()
	f2("lisi",20)
	f3()
	f1()
}

func add(a,b int) (sum int,err error) {
	//多值返回
	if a<0 || b<0 {
		return 0,fmt.Errorf("a and b must be non-negative")
	}
	sum = a+b
	err = nil
	return
}

func sum(params ...int) (sum int) {
	//可变参数(变参数一定要放在参数列表的最后一个)
	for _,i:=range params {
		sum += i
	}

	return
}

func closure() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

//方法-值类型接收者
type Age uint//类型 Age 是方法 String() 的接收者
func (age Age) String() {
	fmt.Printf("%d years\n", age)
}

//方法-指针类型接收者
func (age *Age) Modify(years uint) {
	*age = Age(years)
}

func (age *Age) Add(years uint) {
	*age += Age(years)
}

//声明了一个名为 Person 的结构体类型（Struct Type 数据的集合体）。
type Person struct {
	name string
	age Age
}

//这里如果不是指针类型（p *Person），则p的Modify和Grow方法赋值给变量无法修改结构体的字段值
func (p *Person) Say() {
	fmt.Printf("%s is %d years old\n", p.name, p.age)
}

//方法-指针类型接收者
func (p *Person) Modify(name string, age uint) {
	p.name = name
	p.age = Age(age)
}

func (p *Person) Grow() {
	p.age++
}