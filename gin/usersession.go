package gin

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions/redis"
	redisgo "github.com/gomodule/redigo/redis"
	"github.com/gorilla/sessions"
	"net/http"
)

/**
该扩展包主要用于扩展gin的sessionid没有和用户关联的问题，并且提供对某一用户撤销授权之后，可以删除该用户所有的session
主要方案为：
1. 重写session的Save()方法，在Save（）中增加关联用户信息
2. 撤销用户授权时，通过用户id删除指定session

基于原来的redis存储session的方案，该扩展包采用hash key来保存用户id和该用户对应的sessionid之间的关联关系

用 redis 的 hash key 的 key 存储 userID 标识，用对应的 hash field 存储 sessionID， hash value 的值不重要了，可以不存储内容

主要处理流程为:
1. user login -> generate seesionID and save Session
2. get Session and save user : set userID : sessionID : nil 只存储key，不存储value
这个关联关系的键值对，不关心sessionID的过期情况，但是会出现一个问题，就是如果用户都不是主动过期，就会导致sessionId的key已经过期被清除了
但是这个关联关系的key还没有被清除，即会存在一定的垃圾数据，因此这里做了一个小逻辑来修复这个问题，在用户每次登录的时候，清除掉已经过期的关联关系的key
另外在redis服务层面也有相应的key淘汰策略来处理垃圾数据，所以这个问题不用太过于关注
*/


const sessionPrefix = "session_" //对应 "github.com/gin-contrib/sessions/redis" 中存储sessionid时固定前缀。由于该前缀是在代码中写死的，所以这里只能在代码中写死

//UserSessionExt 提供session的扩展功能，目前提供通过用户ID删除指定用户Session的功能
type UserSessionExt interface {
	DestroyUserSessions(userID string) error
}

//UserSessionStore 用于缓存用户关联数据的store
type UserSessionStore struct {
	*Store               //继承Store，可实现自动顺延功能
	UserKeyPrefix string // 用于保存各业务系统中唯一标示用户信息的key
}

//Save 重写了原来的gin框架自带的RediStore的save方法，在原来的基础上增加了和用户openID关联的方法
func (us *UserSessionStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	//调用store的save方法，延用原来的处理
	err := us.Store.Save(r, w, session)
	if err != nil {
		return err
	}
	return us.saveUser(session)
}

//saveUser 保存session和用户关联的相关信息
func (us *UserSessionStore) saveUser(session *sessions.Session) error {
	//判断当前session是否有userid，有userid则进行关联操作，没有该key则不处理，沿用原来gin session的处理即可
	if v, ok := session.Values[us.UserKeyPrefix].(string); ok && v != "" {
		//关联用户id
		//获取sessionid，这个时候需要将原来的sessionid和用户id建立关联
		//这里用redis的hashmap数据结构来实现，同一个用户的相关关联关系可以用同一个key来表示，对应的field可以用sessionID来标识，value存储当前sessionID的过期时间
		err, rs := redis.GetRedisStore(us.Store)
		if err != nil {
			return err
		}
		conn := rs.Pool.Get()
		if err = conn.Err(); err != nil {
			return err
		}

		//在hset之前，可以先清理一波已经过期的用户信息关联的hash key，以免垃圾数据堆积
		//_, err := conn.Do("HDEL", redisgo.Args{}.Add(us.getUserKey(userID)).AddFlat(sessionIDs)...)
		//主要为了防止如果用户登录很多台设备，又不手动退出，导致有很多已经过期了的hashkey存在，使得key无限增长。

		//清理已经过期的key
		err = us.clearExpiredSessions(conn, v)
		if err != nil {
			return err
		}

		//存储当前用户关联的userID 和 sessionID的hashkey
		//如果当前用户已经登陆过，则续期当前sessionID

		_, err = conn.Do("HSET", us.getUserKey(v), session.ID, nil)
		if err != nil {
			return err
		}

	}
	return nil
}

//clearExpiredSessions 清除当前用户过期的session
func (us *UserSessionStore) clearExpiredSessions(conn redisgo.Conn, userID string) error {
	if conn == nil {
		return errors.New("clearExpiredSessions redis连接为空")
	}
	//遍历当前key下的所有数据，手动判断是否在有效期内
	rpy, err := redisgo.Values(conn.Do("HKEYS", us.getUserKey(userID)))
	if err != nil {
		return err
	}
	fields := make([]string, len(rpy))
	for k, v := range rpy {
		if b, ok := v.([]byte); ok {
			fields[k] = string(b)
		}
	}
	//一次获取所有的sessionid的key的值
	rpy1, err := redisgo.Values(conn.Do("MGET", redisgo.Args{}.AddFlat(fields)...))
	if err != nil && err != redisgo.ErrNil {
		return err
	}

	toClearFields := make([]string, 0)
	//遍历判断当前session是否存在，若存在则不做处理，否则清除该关联关系
	for k, v := range rpy1 {
		if v == nil {
			//已经过期被清除了的sessionid，这里需要对应也删除其对应的userid关联关系的hash key
			toClearFields = append(toClearFields, fields[k])
		}
	}
	return us.deleteUser(conn, userID, toClearFields...)
}

//deleteUser ...
func (us *UserSessionStore) deleteUser(conn redisgo.Conn, userID string, sessionIDs ...string) error {
	if conn == nil {
		return errors.New("deleteUser redis连接为空")
	}

	if len(sessionIDs) <= 0 {
		return nil
	}
	_, err := conn.Do("HDEL", redisgo.Args{}.Add(us.getUserKey(userID)).AddFlat(sessionIDs)...) //
	if err != nil {
		return err
	}
	return nil
}


//getUserKey 获取存储用户session关联关系的field名
func (us *UserSessionStore) getUserKey(userID string) string {
	return fmt.Sprintf("usersession_%s", userID)
}

//DestroyUserSessions 通过用户OpenID删除改用户的所有session
func (us *UserSessionStore) DestroyUserSessions(userID string) error {
	err, rs := redis.GetRedisStore(us.Store)
	if err != nil {
		return err
	}
	conn := rs.Pool.Get()
	if err := conn.Err(); err != nil {
		return err
	}
	//这里只保存关联关系，实际的值并不重要，可以不存储
	rpy, err := redisgo.Values(conn.Do("HKEYS", us.getUserKey(userID)))
	if err != nil {
		return err
	}
	fields := make([]string, len(rpy))
	for k, v := range rpy {
		if b, ok := v.([]byte); ok {
			fields[k] = sessionPrefix + string(b)
		}
	}
	//批量删除对应的session
	_, err = conn.Do("DEL", redisgo.Args{}.AddFlat(fields)...)
	if err != nil {
		return err
	}
	//删除对应的session和关联关系的hash key
	_, err = conn.Do("DEL", redisgo.Args{}.Add(us.getUserKey(userID))) //
	if err != nil {
		return err
	}
	return err
}

//Get 获取session信息
func (us *UserSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(us, name)
}
