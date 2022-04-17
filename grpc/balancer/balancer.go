package balancer

import (
	"fmt"
	"google.golang.org/grpc/balancer"
)

/**
以下完成一套基于hash算法的负载均衡策略
*/

/**
服务启动时，需要注入当前新的balancer组件
注入方式：
balancer.Register(Builder)
 在dial时通过这个option就可以自定义自己的参数了，grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"color"}`)
原理主要是通过该参数改动grpc的ServiceConfig参数：grpc以json的形式提供修改方式，映射的字段为：
```
type jsonSC struct {
	LoadBalancingPolicy *string
	LoadBalancingConfig *internalserviceconfig.BalancerConfig
	MethodConfig        *[]jsonMC
	RetryThrottling     *retryThrottlingPolicy
	HealthCheckConfig   *healthCheckConfig
}
```
然后在resolver 拉起balancer组件就可以通过
switch(balancer)来处理了
 */



type MyBalancer struct {
	cc balancer.ClientConn
}

type MyBalancerBuild struct {
}

func (*MyBalancerBuild) Build(cc balancer.ClientConn, opt balancer.BuildOptions) balancer.Balancer {
	return &MyBalancer{cc: cc}
}

func (*MyBalancerBuild) Name() string {
	return "mybalancer"
}

func (mp *MyBalancer) UpdateClientConnState(state balancer.ClientConnState) error {
	//更新连接状态，这个函数是全局中最重要的一个函数
	//也就是靠这个函数来更新管理整个服务连接状态的
	return nil
}

func (mp *MyBalancer) ResolverError(err error) {
	//nothing todo
	fmt.Println(err)
}

func (mp *MyBalancer) UpdateSubConnState(conn balancer.SubConn, state balancer.SubConnState) {
	//更新子链接状态，这个是由每个子连接监听器做管理的
	if conn == nil {
		return
	}

}

func (mp *MyBalancer) Close() {
	//优雅地关闭负载均衡连接器，需要释放所有的连接

}

