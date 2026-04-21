package main

import (
	"fmt"
	"sort"
	"unicode/utf8"
)

/**
 * 集合类型：数组、切片、Map
 */
func main() {
	//数组：[5]string 和 [4]string 不是同一种类型，也就是说长度也是数组类型的一部分
	array:=[5]string{"a","b","c","d","e"}
	fmt.Println(array[2])

	array1 := [...]int{1,2,3,4,5}
	fmt.Println(array1)

	array2 := [5]string{1:"a", 3:"b"}
	fmt.Println(array2)

	for i:=0; i<5; i++{
		fmt.Printf("数组索引:%d,对应值:%s\n", i, array2[i])
	}

	for i,v:=range array2{
		fmt.Printf("数组索引:%d,对应值:%s\n", i, v)
	}

	for _,v:=range array2{
		fmt.Printf("对应值:%s\n", v)
	}

	array3 := [2][3]string{{"a","b","c"},{"d","e","f"}}
	for k1,v1 :=range array3{
		fmt.Printf("第一层：%d\n",k1)
		for k2,v2 :=range v1{
			fmt.Printf("第二层：%d,对应值:%s\n",k2,v2)
		}
	}

	//Slice（切片）：切片是基于数组实现的，它的底层就是一个数组。对数组任意分隔，就可以得到一个切片
	slice:=array[2:4]//基于数组生成切片，包含索引start，但是不包含索引end
	fmt.Println(slice)

	//经过切片后，切片的索引范围改变了
	for i,v:=range slice{
		fmt.Printf("切片索引:%d,对应值:%s\n", i, v)
	}

	//切片表达式 array[start:end] 中的 start 和 end 索引都是可以省略的
	slice1:=array[:4]
	slice2:=array[2:]
	slice3:=array[:]
	fmt.Println(slice1,slice2,slice3)

	//切片修改,一旦修改切片的元素值，那么底层数组对应的值也会被修改
	slice[1] = "f"
	fmt.Println(slice, array)

	//切片声明
	slice4 := make([]string,4,8)//切片的容量不能比切片的长度小,其他的内存空间处于空闲状态(想要高性能，预先分配容量)
	fmt.Println(len(slice4),cap(slice4))

	slice5 := []string{"a","b","c","d","e"}
	fmt.Println(len(slice5),cap(slice5))

	//在创建新切片的时候，最好要让新切片的长度和容量一样
	//追加一个元素
	slice6 := append(slice5, "f")
	//追加多个元素
	slice7 := append(slice5, "f", "g")
	//追加另一个切片
	slice8 := append(slice5, slice...)
	slice8 = append(slice8, slice1...)
	slice9 := append(append(slice5, slice...), slice2...)
	slice10 := mergeSlices(slice,slice1,slice2)
	fmt.Println(slice6,slice7,slice8,slice9,slice10)

	// 二维切片：2 行 3 列
	arr := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}

	fmt.Println(arr) // [[1 2 3] [4 5 6]]

	// 遍历
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			fmt.Printf("arr[%d][%d] = %d\n", i, j, arr[i][j])
		}
	}

	// 先创建 2 行，每行再 make 3 个元素
	arr1 := make([][]int, 2)
	for i := range arr1 {
		arr1[i] = make([]int, 3)
	}

	// 赋值
	arr1[0][1] = 100
	arr1[1][2] = 200
	fmt.Println(arr1)

	//Map（映射）
	nameAgeMap:=make(map[string]int)
	nameAgeMap["zhangsan"] = 20
	nameAgeMap["lisi"] = 21
	nameAgeMap["wangwu"] = 22
	fmt.Println(nameAgeMap)

	//Map的声明
	nameAgeMap1:=map[string]int{"zhangsan":20,"lisi":21,"wangwu":22}
	fmt.Println(nameAgeMap1)

	ageZhangsan:=nameAgeMap1["zhangsan"]
	fmt.Println(ageZhangsan)

	//map 可以获取不存在的 K-V 键值对，如果 Key 不存在，返回的 Value 是该类型的零值
	age,ok:=nameAgeMap1["zhangsan"]
	if ok{
		fmt.Println(age)
	}else{
		fmt.Println("zhangsan not found")
	}

	delete(nameAgeMap1,"zhangsan")

	//map 的遍历是无序的(如果想按顺序遍历，可以先获取所有的 Key，并对 Key 排序，然后根据排序好的 Key 获取对应的 Value)
	for k,v:=range nameAgeMap1{
		fmt.Printf("key:%s,value:%d\n", k, v)
	}

	//顺序打印map
	user := map[string]int{
		"bob":   20,
		"alice": 18,
		"tom":   22,
	}

	//把所有 key 拿出来
	var keys []string
	for k := range user {
		keys = append(keys, k)
	}

	//对 key 排序(如果 key 是 int 类型,把 sort.Strings 换成 sort.Ints 就行)
	sort.Strings(keys)

	fmt.Println("=== 顺序打印 map ===")
	for _, k := range keys {
		fmt.Printf("key: %s, value: %d\n", k, user[k])
	}


	fmt.Println(len(user))

	// 二维 map
	m := make(map[string]map[string]int)

	// 第一层 key
	m["user1"] = make(map[string]int)
	// 第二层 key
	m["user1"]["age"] = 20
	m["user1"]["score"] = 99

	m["user2"] = make(map[string]int)
	m["user2"]["age"] = 22
	m["user2"]["score"] = 88

	fmt.Println(m)
	//安全判断
	if val, ok := m["user1"]; ok {
		fmt.Println(val["age"])
	}

	//简化写法（直接初始化）
	m1 := map[string]map[string]int{
		"a": {
			"x": 1,
			"y": 2,
		},
		"b": {
			"x": 3,
			"y": 4,
		},
	}

	for k1,v1 :=range m1{
		fmt.Println("第一层 key:", k1)
		for k2,v2 :=range v1{
			fmt.Printf("  %s: %d\n", k2, v2)
		}
	}

	//String 和 []byte
	s := "hello, world 你好，世界"//UTF8 编码下，一个汉字对应三个字节
	fmt.Println(len(s))
	fmt.Println(utf8.RuneCountInString(s))//把一个汉字当成一个长度计算

	bs := []byte(s)//转为字节切片 []byte
	fmt.Println(bs)
	fmt.Println(s[0],s[1],s[15])

	//for range 循环在处理字符串的时候，自动地隐式解码 unicode 字符串
	for i,r :=range s{
		fmt.Printf("%d:%c\n", i, r)
	}
}

/**
合并切片
*/
func mergeSlices(slices ...[]string) []string{
	var res []string
	for _,s :=range slices{
		res = append(res, s...)
	}

	return res
}