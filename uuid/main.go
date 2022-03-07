package main

import (
	"fmt"
)

type Name struct {
	AA string
}
func main() {

	s := make([]*Name, 2)
	for k, _ := range s {
		if s[k] == nil {
			fmt.Println("nil")
		}
	}

}
