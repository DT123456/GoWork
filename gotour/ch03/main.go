package main

import "fmt"

/**
 * 控制结构：if、for、switch 逻辑语句
 */
func main() {
	i := 10
	if i > 0 {
		fmt.Println("i is greater than 0")
	} else {
		fmt.Println("i is less than 0")
	}

    if i1:=6;i1 >10 {
        fmt.Println("i1>10")
    } else if  i1>5 && i1<=10 {
        fmt.Println("5<i1<=10")
    } else {
        fmt.Println("i1<=5")
    }

	switch i2:=6;{
		case i2>10:
			fmt.Println("i2>10")
		case i2>5 && i2<=10:
			fmt.Println("5<i2<=10")
		default:
			fmt.Println("i2<=5")
	}

	switch j:=1;j {
		case 1:
			fallthrough
		case 2:
			fmt.Println("1")
		default:
			fmt.Println("没有匹配")
	}

	switch 2>1 {
		case true:
			fmt.Println("2>1")
		case false:
			fmt.Println("2<=1")
	}

	sum:=0
	for i:=1;i<=100;i++ {
		sum+=i
	}
	fmt.Println("the sum is",sum)

	sum1:=0
	i1:=1
	for {
		sum1+=i1
		i1++
		if i1>100 {
			break
		}
	}
	fmt.Println("the sum1 is",sum1)

	sum = 0
	for i:=1; i<100; i++{
	if i%2!=0 {
		continue
	}
	sum+=i
	}
	fmt.Println("the sum is",sum)
}