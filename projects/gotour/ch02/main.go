package main

import (
	"fmt"
	"strconv"
	"strings"
)

/**
 * 数据类型
 */
func main() {
	var i int = 10
	fmt.Println(i)

	var (
		j int = 11
		k int = 21
	)
	fmt.Println(j, k)

	var f32 float32 = 1.1
	var f64 float64 = 1.23456
	fmt.Println("f32 is",f32,",f64 is",f64)

	var bf bool = true
	var bf2 bool = false
	fmt.Println(bf,bf2)

	var s1 string = "hello"
	var s2 string = "world"
	fmt.Println("s1 is",s1,",s2 is",s2)
	fmt.Println(s1+s2) //s1+=s2 同理

	var zi int
	var zf float64
	var zb bool
	var zs string
	fmt.Println(zi,zf,zb,zs)

	i1 := 10
	i2 := 20
	bool1 := true
	s3 := "hello"
	fmt.Println(i1,i2,bool1,s3)

	pi := &i1
	fmt.Println(*pi)

	i1 = 20
	fmt.Println("i的新值是",i1,*pi)

	const name = "zhangsan"
	const (
		one = 1
		two = 2
		three = 3
		four = 4
	)
	const (
		one1 = iota + 1
		two1
		three1
		four1
	)
	fmt.Println(name)
	fmt.Println(one,two,three,four)
	fmt.Println(one1,two1,three1,four1)

	i2s := strconv.Itoa(i1)
	s2i,err := strconv.Atoi(i2s)
	fmt.Println(i2s,s2i,err)

	i2f := float64(i1)
	f2i := int(i2f)
	fmt.Println(i2f,f2i)

	fmt.Println(strings.HasSuffix(s1,"o")) //s1是否以O结尾
	fmt.Println(strings.HasPrefix(s1,"h")) //s1是否以H开头
	fmt.Println(strings.Contains(s1,"e")) //s1是否包含e
	fmt.Println(strings.Index(s1,"e")) //s1中e的索引
	fmt.Println(strings.Count(s1,"l")) //s1中l的个数
	fmt.Println(strings.Replace(s1,"l","L",1)) //s1中将l替换为L，只替换一个
	fmt.Println(strings.ReplaceAll(s1,"l","L")) //s1中将l替换为L，全部替换
	fmt.Println(strings.ToLower(s1)) //s1转换为小写
	fmt.Println(strings.ToUpper(s1)) //s1转换为大写
	fmt.Println(strings.TrimSpace(s1)) //s1去除首尾空格
	fmt.Println(strings.TrimLeft(s1,"h")) //s1去除左侧H
	fmt.Println(strings.TrimRight(s1,"o")) //s1去除右侧O
	fmt.Println(strings.TrimPrefix(s1,"h")) //s1去除前缀H
	fmt.Println(strings.TrimSuffix(s1,"o")) //s1去除后缀O
	fmt.Println(strings.Split(s1,"e")) //s1按e分割
	fmt.Println(strings.Join(strings.Split(s1,"e"),"E")) //s1按e分割，并将e替换为E
	fmt.Println(strings.Repeat(s1,2)) //s1重复2次
	
}