package main
//
//
//
//
///**-------- 不用泛型-------**/
//type GenericType interface{}
//
//type GenericListType interface {}
//
//// List ...
//type List interface {
//	GetAll() GenericListType //获取所有类型
//	Add(GenericType)       //添加一个元素
//	Sum() GenericType      //求整个集合元素的值
//}
//
////GenericList ...
//type GenericList struct {
//	List GenericListType
//}
//
//func (gl *GenericList) GetAll() GenericListType {
//	return gl.List
//}
//
//func (gl *GenericList) Add(item GenericType) {
//
//	switch gl.List.(type) {
//	case []int:
//		gl.List = append(gl.List.([]int), item.(int))
//	case []float32:
//		gl.List = append(gl.List.([]float32), item.(float32))
//	case []string:
//		gl.List = append(gl.List.([]string), item.(string))
//	default:
//		panic("参数类型异常")
//	}
//}
//
//func (gl *GenericList) Sum() GenericType {
//
//	switch gl.List.(type) {
//	case []int:
//		//逐个进行类型检查
//		var count int
//		for _, v := range gl.List.([]int) {
//			count = count + v
//		}
//		return count
//	case []float32:
//		var count float32
//		for _, v := range gl.List.([]float32) {
//			count = count + v
//		}
//		return count
//	case []string:
//		var count string
//		for _, v := range gl.List.([]string) {
//			count = count + v
//		}
//		return count
//	default:
//		panic("参数类型异常")
//
//	}
//
//}
//
//
