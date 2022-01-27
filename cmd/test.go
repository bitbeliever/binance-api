package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {

	var a st
	fmt.Println(a)
	a.chang()
	fmt.Println(a)

	resp, err := http.Get("http://www.dookbook.com/books/booksMore/id-3?&p=6")
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(b))
}

type st struct {
	a int
}

func (s *st) chang() {
	s.a = 1
}
