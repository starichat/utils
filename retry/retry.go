package utils

import (
	"time"
)

/**
一个简单的重试工具类
*/

//HandlerFunc ...
type HandlerFunc func() error

//Config 重试相关配置
type Config struct {
	maxRetryCount  int           //最大重试次数
	delayTime      time.Duration //延时多久重试
	maxHandlerTime time.Duration //整个请求最大请求时间
}

func newDefaultConfig() *Config {
	return &Config{
		maxRetryCount:  3,
		delayTime:      10 * time.Second,
		maxHandlerTime: 5 * time.Minute,
	}
}

//Option ...
type Option func(config *Config)

//Do 执行函数
func Do(f HandlerFunc, opts ...Option) error {
	//解析配置
	conf := newDefaultConfig()
	for _, opt := range opts {
		opt(conf)
	}
	var lastErr error
	//并且延时delay时间后再次重试
	lastErrIndex := 0
	for lastErrIndex < conf.maxRetryCount {
		if err := f(); err != nil {
			lastErrIndex++
			//判断是否是最后一次重试,如果是则直接返回错误,否则延时重试
			if lastErrIndex == conf.maxRetryCount {
				return err
			}
			//延时重试
			select {
			case <-time.After(conf.delayTime):
				continue
			}

		} else {
			return nil
		}

	}
	//监控整个请求的最大请求时间,(包括所有重试请求的时间),避免请求一直阻塞,超时直接报错.todo
	return lastErr
}

//WithRetryCount 增加最大重试次数
func WithRetryCount(count int) Option {
	return func(config *Config) {
		config.maxRetryCount = count
	}
}

//WithDelayTime 配置延时时间
func WithDelayTime(delay time.Duration) Option {
	return func(config *Config) {
		config.delayTime = delay
	}
}
