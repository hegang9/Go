package main // 必须声明当前 go 文件属于哪个包，入口文件必须声明为 main 包

import (
	"fmt"
)

func main() {
	var a, b int
	fmt.Scanln(&a, &b)

	if a > b {
		fmt.Println("a is greater than b")
	} else {
		fmt.Println("a is less than or equal to b")
	}

	switch {
	case a > b:
		fmt.Println("a is greater than b")
	case a < b:
		fmt.Println("a is less than b")
	case a == b:
		fmt.Println("a is equal to b")
	}

	for ; a < b; a++ {
		fmt.Println("a is less than b, incrementing a:", a)
	}

	nums := []int{10, 20, 30}
	for index, value := range nums {
		fmt.Printf("索引: %d, 值: %d\n", index, value)
	}

	m := make(map[string]int)
	m["one"] = 1
	fmt.Println("Map value for 'one':", m["one"])
}
