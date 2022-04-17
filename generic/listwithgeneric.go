package main



/**-------- 使用泛型-------**/
type GenericType interface {
	int | float32 | float64 | string
}

//type GenericListType []GenericType


// List ...
type List[V GenericType] interface {
	GetAll() []V //获取所有类型
	Add(V)       //添加一个元素
	Sum() V      //求整个集合元素的值
}

//GenericList ...
type GenericList[V GenericType] struct {
	List []V
}

func (gl *GenericList[V]) GetAll() []V {
	return gl.List
}

func (gl *GenericList[V]) Add(item V) {
	gl.List = append(gl.List, item)
}

func (gl *GenericList[V]) Sum() V {
	var res V
	for _, v := range gl.List {
		res = res + v
	}
	return res

}


