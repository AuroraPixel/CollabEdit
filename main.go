package main

import (
	"fmt"
)

type Person interface {
}

type User struct {
	Name string
	Age  int
}

func main() {

	var intMap = make(map[int]int)
	intMap[1] = 1
	intMap[2] = 2
	fmt.Println(intMap[55])
}
