package main

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

/**
 * 运行时反射：字符串和结构体之间如何转换？
 */
func main() {
	//Println 参数a ...interface{}，
	// ...interface{} 表示可变参数，以传递多个任意类型的参数
	//interface{} 是空接口，可以表示任何类型
	fmt.Println("Hello, World!")

	fmt.Println("===reflect.Value 和 reflect.Type===")
	i:=3
	// 通过 reflect.ValueOf 函数把任意类型的对象转为一个 reflect.Value 类型的对象，
	// reflect.Value 是一个结构体类型，包含了一个 interface{} 类型的字段，存储了原始对象的值。
	// reflect.TypeOf同理
	//int to reflect.Value
	iv:=reflect.ValueOf(i)
	it:=reflect.TypeOf(i)
	fmt.Println(iv,it)//3 int

	// 逆向转回来
	//reflect.Value to int
	i1:=iv.Interface().(int)
	fmt.Println(i1)

	// 修改对应的值
	ipv:=reflect.ValueOf(&i) // reflect.ValueOf 函数返回的是一份值的拷贝，传入变量的指针才可以修改
	ipv.Elem().SetInt(4) // Elem() 方法返回一个新的 reflect.Value，表示指针指向的值。SetInt() 方法将该值设置为 4。
	fmt.Println(i)

	p:=person{Name: "飞雪无情",Age: 20}
	ppv:=reflect.ValueOf(&p)
	ppv.Elem().Field(0).SetString("张三") // Field(0) 方法返回一个新的 reflect.Value，表示结构体的第一个字段（Name）。SetString() 方法将该字段设置为 "张三"。
	ppv.Elem().Field(1).SetInt(30) // Field(1) 方法返回一个新的 reflect.Value，表示结构体的第二个字段（Age）。SetInt() 方法将该字段设置为 30。
	fmt.Println(p)
	fmt.Println(ppv.Kind())
	pv:=reflect.ValueOf(p)
	fmt.Println(pv.Kind())

	pt:=reflect.TypeOf(p)
	//遍历person的字段
	for i:=0;i<pt.NumField();i++{ // NumField() 方法返回结构体的字段数量，Field(i) 方法返回一个 reflect.StructField 类型的对象，包含了字段的名称、类型、标签等信息。
		fmt.Println("字段：",pt.Field(i).Name)
	}
	//遍历person的方法
	for i:=0;i<pt.NumMethod();i++{// NumMethod() 方法返回结构体的方法数量，Method(i) 方法返回一个 reflect.Method 类型的对象，包含了方法的名称、类型、函数等信息。
		fmt.Println("方法：",pt.Method(i).Name)
	}

	// 判断是否实现了某个接口
	stringerType:=reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	writerType:=reflect.TypeOf((*io.Writer)(nil)).Elem()
	fmt.Println("是否实现了fmt.Stringer：",pt.Implements(stringerType))
	fmt.Println("是否实现了io.Writer：",pt.Implements(writerType))

	fmt.Println("===字符串和结构体互转===")
	//JSON 和 Struct 互转
	//struct to json
	jsonB,err:=json.Marshal(p)
	if err==nil {
		fmt.Println(string(jsonB))
	}
	//json to struct
	respJSON:="{\"Name\":\"李四\",\"Age\":40}"
	json.Unmarshal([]byte(respJSON),&p)
	fmt.Println(p)

	respJSON1:="{\"name\":\"李四\",\"age\":40}"
	json.Unmarshal([]byte(respJSON1),&p)
	fmt.Println(p)

	//遍历person字段中key为json的tag
	for i:=0;i<pt.NumField();i++{
		sf:=pt.Field(i)
		//可以多个tag
		fmt.Printf("字段%s上,json tag为%s\n",sf.Name,sf.Tag.Get("json"))
		fmt.Printf("字段%s上,bson tag为%s\n",sf.Name,sf.Tag.Get("bson"))
	}

	//自己实现的struct to json
	jsonBuilder:=strings.Builder{}
	jsonBuilder.WriteString("{")
	num:=pt.NumField()
	for i:=0;i<num;i++{
		jsonTag:=pt.Field(i).Tag.Get("json") //获取json tag
		jsonBuilder.WriteString("\""+jsonTag+"\"")
		jsonBuilder.WriteString(":")
		//获取字段的值
		jsonBuilder.WriteString(fmt.Sprintf("\"%v\"",pv.Field(i)))
		if i<num-1{
			jsonBuilder.WriteString(",")
		}
	}
	jsonBuilder.WriteString("}")
	fmt.Println(jsonBuilder.String())//打印json字符串

	fmt.Println("===反射的三大定律===")
	// 反射的三大定律：
	// 1. 反射可以将 interface{} 转换为具体类型
	// 2. 反射可以将具体类型转换为 interface{}
	// 3. 如果修改反射对象，值必须是可设置的（settable）

	fmt.Println("===使用反射调用方法===")
	//反射调用person的Print方法
	mPrint:=pv.MethodByName("Print")
	args:=[]reflect.Value{reflect.ValueOf("登录")} //声明参数，类型是 []reflect.Value（切片） 调用Print方法，传入参数"登录"
	mPrint.Call(args)
}

type person struct {
	Name string `json:"name"  bson:"b_name"` //添加 tag: 结构体标签，指定字段在 JSON 中的名称
	Age int `json:"age" bson:"b_age"`
}

func (p person) String() string{
   return fmt.Sprintf("Name is %s,Age is %d",p.Name,p.Age)
}

func (p person) Print(prefix string){
   fmt.Printf("%s:Name is %s,Age is %d\n",prefix,p.Name,p.Age)
}


