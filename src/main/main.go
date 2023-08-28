package main

import (
	"fmt"
)

type A struct {
	num int
}

func (a A) f() {
	fmt.Println(a.num)
}

func main() {
	a := A{num: 10}
	a.f()
}
