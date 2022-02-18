package main

import (
	"errors"
	"fmt"
)

/**
热冷库操作
 */

/**
1. 迁移数据库数据到新库
2.
 */

func MigrateData() {

}

func main() {
	s := make([]error,0)
	err := errors.New("1")
	s = append(s, err)
	err = errors.New("2")
	s = append(s, err)
	for _, v := range s {
		fmt.Printf("当前元素地址:%p, %v\n",v, v)
	}
	fmt.Println(s)
}
