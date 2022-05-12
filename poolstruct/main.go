package main

import "container/ring"

type Pool struct {
	idle *ring.Ring
	active *ring.Ring
}

type Item struct {
	id int
	active bool
}

func main() {

}
