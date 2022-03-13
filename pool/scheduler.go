package pool


type scheduler struct {
	strategy Strategy //调度策略
	schedule Schedule //调度接口

}

//Schedule 调度器接口
type Schedule interface {
	Put() (bool, error) //添加连接
	Get() (*Conn, error) //获取连接
	Destroy() error //关闭连接

}

type strategy struct {
	name string
}

//Strategy 调度策略
type Strategy interface {
	Build() strategy
	Name() string
}



