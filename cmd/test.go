package main

import (
	"fmt"
)

func main() {

	var a st
	fmt.Println(a)
	a.chang()
	fmt.Println(a)
}

type st struct {
	a int
}

func (s *st) chang() {
	s.a = 1
}
