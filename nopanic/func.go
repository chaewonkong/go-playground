package main

import "fmt"

type NilInterface interface {
	Hello()
}

type ParentStruct struct {
	Src NilInterface
}

var _ NilInterface = &ParentStruct{}

func NewParentStruct(child NilInterface) NilInterface {
	return &ParentStruct{
		Src: child,
	}
}

func (p *ParentStruct) Hello() {
	p.Src.Hello()
}

// test
type ChildStruct struct{}

func NewChildStruct() *ChildStruct {
	return nil
}

func (c *ChildStruct) Hello() {
	fmt.Println("hello")
}
