package main

import (
	"fmt"
	"unsafe"
)

/**
 * 非类型安全：unsafe
 */
func main() {
	fmt.Println("===指针类型转换===")
	// Go 语言是强类型，不允许两个指针类型进行转换的。
	// 但是 Go 语言提供了 unsafe 包，可以绕过类型安全的限制，进行指针类型转换。
	// unsafe.Pointer 是一个特殊的指针类型，可以表示任何类型的指针。
	// uintptr 是一个无符号整数类型，可以表示指针的地址。

	fmt.Println("===unsafe.Pointer 类型转换===")
	var a int = 10
	var intP *int = &a
	fmt.Println("intP变量的值为:",intP)
	fmt.Println("intP变量的值为:",*intP)

	// 将 intP 转换为 unsafe.Pointer 类型
	unsafePtr := unsafe.Pointer(intP)
	fmt.Println("unsafePtr变量的值为:",unsafePtr)

	var fp *float64 = (*float64)(unsafePtr)
	*fp = *fp * 3
	fmt.Println("a变量的值为:",a)

	// 将 unsafe.Pointer 转换回 *int 类型
	intP1 := (*int)(unsafePtr)
	fmt.Println("intP1变量的值为:",intP1)
	fmt.Println("intP1变量的值为:",*intP1)

	fmt.Println("===uintptr 的转换===")
	// 将 unsafe.Pointer 转换为 uintptr 类型 uintptr 可以进行算术运算，unsafe.Pointer 不可以
	uintptrVal := uintptr(unsafePtr)
	fmt.Println("uintptrVal变量的值为:",uintptrVal)

	// 将 uintptr 转换回 unsafe.Pointer 类型
	unsafePtr2 := unsafe.Pointer(uintptrVal)
	fmt.Println("unsafePtr2变量的值为:",unsafePtr2)

	p := new(person)
    pName:=(*string)(unsafe.Pointer(p))
   	*pName = "飞雪无情"
	//Age并不是person的第一个字段，所以需要进行偏移，这样才能正确定位到Age字段这块内存，才可以正确的修改
	// unsafe.Offsetof(p.Age) 返回 person 结构体中 Age 字段相对于结构体起始地址的偏移量，单位是字节。通过将 p 的地址加上这个偏移量，就可以得到 Age 字段的地址。
	// 指针运算完毕后，通过 unsafe.Pointer 转换为真实的指针类型,进行赋值或取值操作
	pAge := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Offsetof(p.Age)))
	*pAge = 20
	fmt.Println("person变量的值为:",p)
	fmt.Println("person变量的值为:",*pName,*pAge)

	fmt.Println("===指针转换规则===")
	// 任何类型的 *T 都可以转换为 unsafe.Pointer；
	// unsafe.Pointer 也可以转换为任何类型的 *T；
	// unsafe.Pointer 可以转换为 uintptr；
	// uintptr 也可以转换为 unsafe.Pointer。

	fmt.Println("===指针类型转换的风险===")
	// 指针类型转换的风险在于，如果转换后的指针类型与原始指针类型不兼容，可能会导致程序崩溃或者数据损坏。
	// 因此，在进行指针类型转换时，必须确保转换后的指针类型与原始指针类型兼容，并且要非常小心地使用 unsafe 包。	
	
	fmt.Println("===unsafe.Sizeof===")
	//返回一个类型所占用的内存大小，大小只与类型有关，和类型对应的变量存储的内容大小无关
	fmt.Println("a的大小为:", unsafe.Sizeof(a))
	fmt.Println("person结构体的大小为:", unsafe.Sizeof(p))
	fmt.Println("person结构体的大小为:", unsafe.Sizeof(*p))

	fmt.Println(unsafe.Sizeof(true))
	fmt.Println(unsafe.Sizeof(int8(0)))
	fmt.Println(unsafe.Sizeof(int16(10)))
	fmt.Println(unsafe.Sizeof(int32(10000000)))
	fmt.Println(unsafe.Sizeof(int64(10000000000000)))
	fmt.Println(unsafe.Sizeof(int(10000000000000000)))
	fmt.Println(unsafe.Sizeof(string("飞雪无情")))
	fmt.Println(unsafe.Sizeof([]string{"飞雪u无情","张三"}))
}

type person struct {
	Name string
	Age int
}