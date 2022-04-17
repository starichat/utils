package main

import "fmt"

//
////定义支持泛型的函数
//func compare[T comparable](a,b T) T {
//	if a != b {
//		return a
//	}
//	return b
//}
//
//func min[T int64|float32](a, b T) T {
//if a < b {
//return a
//}
//return b
//}


type Obj[T Numerber] struct {}

func (o *Obj[T]) min(a, b int) int {
	if a < b {
		return a
	}
	return b

}

type Numerber interface {
	int64 | float32 | float64
}

type ObjInterface[T int64 | float32| float64 ] interface {
	Min(a,b T) T
	min(a,b int) int
}

func (o *Obj[T]) Min(a,b T) T {
	if a < b {
		return a
	}
	return b
}


//type OOO interface {
//	Min[T int64 | float32](a,b T) T
//	min(a,b int) int
//}

func NewII[T int64|float32|float64](obj *Obj[T]) ObjInterface[T] {
	return obj
}


func main() {
	//
	//fmt.Println(compare[int](1,2))
	//fmt.Println(compare[string]("he","ze"))
	//fmt.Println(compare[float32](1.222,1.3333))
	//fmt.Println(compare[float64](1.222,1.3333))
	//
	//min[int64](1,2)
	o := Obj[int64]{}
	fmt.Println(o.min(1,2))
	o1 := Obj[float32]{}
	fmt.Println(o1.min(1.0,2.0))
	NewII[int64](&o)

	l0 :=
	l1 :=GenericList[float32]{List:  make([]float32,0)}
	l2 :=
	}



	ml := []GenericList{
		GenericList[int]{List: make([]int,0)},
		GenericList[float32]{List:  make([]float32,0)},
		GenericList[string]{List: make([]string,0),
	}
	ml[0].Add(1)
	ml[0].Add(2)
	ml[0].Add(3)
	ml[0].Add(4)
	ml[0].Add(5)
	fmt.Printf("[int]元素的值:%+v\n", l0.GetAll())
	fmt.Printf("[int]集合求和:%v\n", l0.Sum())

	ml[1].Add(1.5)
	ml[1].Add(2.5)
	ml[1].Add(3.6)
	ml[1].Add(4.7)
	ml[1].Add(5.8)
	fmt.Printf("[float32]元素的值:%+v\n", l1.GetAll())
	fmt.Printf("[float32]集合求和:%v\n", l1.Sum())


	ml[2].Add("hello")
	ml[2].Add(" go")
	ml[2].Add(" generic")
	ml[2].Add(" coming")
	ml[2].Add(" yeah!!!")
	fmt.Printf("[string]元素的值:%+v\n", l2.GetAll())
	fmt.Printf("[string]集合求和:%v\n", l2.Sum())
}