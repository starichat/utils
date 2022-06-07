package algorithm

import "container/ring"

type RingPool struct {
	Idle *ring.Ring
	Active *ring.Ring
}

type Item struct {
	ID string
	Status int
}

type ListPool struct {
	Idle []*Item
}



type Pool interface {
	Get() Item
	Put(i *Item)
	Update()
}

func (l ListPool) Get() Item {
	//随机获取数据，从idle中获取数据
	if getIndex()
	if v := l.Idle[getIndex()];  {

	}
}

func getIndex() int {
	return 1
}

func (l ListPool) Put(i *Item) {
	panic("implement me")
}

func (l ListPool) Update() {
	panic("implement me")
}


