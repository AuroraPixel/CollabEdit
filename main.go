package main

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type Person interface {
	Write() string
	TODO() string
}

type Student struct {
	Person
	Name string
	Age  int
}

func (u Student) Write() string {
	return "写作业"
}

func (u Student) TODO() string {
	return "上学"
}

type Worker struct {
	Person
	Name string
	Age  int
}

func (u Worker) Write() string {
	return "写PPT"
}

func (u Worker) TODO() string {
	return "上班"
}

func main() {
	p := &Worker{}

	p1 := interfaceToPersion(p)
	fmt.Println(p1.TODO())

	s := &Student{}
	s1 := interfaceToPersion(s)
	fmt.Println(s1.TODO())

}

func interfaceToPersion(i interface{}) Person {

	switch i.(type) {
	case Person:
		fmt.Println("Person")
		break

	case *Person:
		fmt.Println("*Person")
		break
	case Worker:
		fmt.Println("Worker")
		break
	case *Worker:
		fmt.Println("*Worker")
		break

	case Student:
		fmt.Println("Student")
		break
	case *Student:
		fmt.Println("*Student")
		break

	default:
		fmt.Println("default")
		break

	}

	var p Person
	mapstructure.Decode(i, &p)
	fmt.Println(p.TODO())
	return p
}
