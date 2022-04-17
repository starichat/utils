package gin

import (
	"encoding/base32"
	"net/http"
	"strings"
	"time"

	gsessions "github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

/*
 * http中的回话管理使用方式多种多样，但是最有效的还是基于cookie的实现方案，
 * 在没有特殊情况，建议使用基于cookie的解决方案。
 *
 * 基于cookie的方式进行回话管理更具cookie的特性和业务场景，可以分解为：
 * session
 * ├── 有指定有效时间的回话(例如两小时)
 * │   ├── 有效期从颁发开始计算，时间固定
 * │   ├── 有效期从有数据访问进行顺延
 * ├── 没有固定有效期(浏览器关闭会话结束(隐藏有固定时效))
 *
 * 现有session解决方案中，基本使用redis作为存储，现基于gin矿建的session中间件中redis存储进行拓展
 */

const (
	sessionNoChange   = "__sessionNoChange" //用于标记sess是否进行了set操作
	sessionConstStyle = "__sessionConst"    //用于写入session时间
)

//Store 基于gin的session中间件进行的扩展
//Session 默认是调用Save方法则将有效期顺延指定时间，不调用Save及时有请求也不进行顺延
type Store struct {
	pairs       []securecookie.Codec //用于对cookie数据进行加密
	redis.Store                      //redis store
	MaxAgeZero  bool                 //是否是基于浏览器的(cookie没有有效期)
	ConstExp    bool                 //固定有效期(由于cookie的单位是秒，会产生秒级别的波动)
	AutoRefresh bool                 //自动进行顺延(及时不调用save方法)   需要重写
}

//Save 保存session
func (s *Store) Save(r *http.Request, w http.ResponseWriter, sess *sessions.Session) error {
	if c, err := r.Cookie(sessionConstStyle); s.ConstExp && err == http.ErrNoCookie {
		//写入cookie生成时间
		t, _ := time.Now().MarshalJSON()
		opt := &sessions.Options{}
		*opt = *sess.Options
		http.SetCookie(w, sessions.NewCookie(sessionConstStyle, strings.Trim(string(t), "\""), opt)) //写入时间
	} else if s.ConstExp && err == nil { //存在session
		t := time.Now()
		err := (&t).UnmarshalJSON([]byte("\"" + c.Value + "\""))
		if err != nil {
			return err
		}
		sess.Options.MaxAge -= int(time.Now().Sub(t).Seconds()) //剩余时间
	}

	if _, has := sess.Values[sessionNoChange]; has && sess.ID != "" { //必须是存在sessionid 且明确标记为不进行值更新操作的
		delete(sess.Values, sessionNoChange)
		return s.noChangeSave(r, w, sess)
	}

	//调用store的save方法等价于AutoRefresh
	err := s.Store.Save(r, w, sess) //调用save方法，肯定会重写cookie信息，需要要再次写入cookie信息进行数据覆盖
	if err != nil {
		return err
	}
	return s.rewriteCookie(r, w, sess)
}

func (s *Store) noChangeSave(r *http.Request, w http.ResponseWriter, sess *sessions.Session) error {
	//先进行续约
	if sess.ID == "" {
		sess.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
	}
	err, rs := redis.GetRedisStore(s.Store)
	if err != nil {
		return err
	}
	conn := rs.Pool.Get()
	defer conn.Close()
	age := sess.Options.MaxAge
	if age == 0 {
		age = rs.DefaultMaxAge
	}
	_, err = conn.Do("EXPIRE", "session_"+sess.ID, age) //
	if err != nil {
		return err
	}

	if s.MaxAgeZero {
		return s.rewriteCookie(r, w, sess)
	}

	encoded, err := securecookie.EncodeMulti(sess.Name(), sess.ID, s.pairs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, sessions.NewCookie(sess.Name(), encoded, sess.Options))
	return nil
}

func (s *Store) rewriteCookie(r *http.Request, w http.ResponseWriter, sess *sessions.Session) error {
	if s.MaxAgeZero {
		//cookie的有效期设置为0
		opt := &sessions.Options{}
		*opt = *sess.Options
		opt.MaxAge = 0
		encoded, _ := securecookie.EncodeMulti(sess.Name(), sess.ID, s.pairs...)
		http.SetCookie(w, sessions.NewCookie(sess.Name(), encoded, opt))
		return nil
	}
	return nil
}

//Get 获取session信息
func (s *Store) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

//MarkSessionNoChange 用于标记session中数据未进行改动，但需要调用Save方法进行数据有效期延迟
func MarkSessionNoChange(sess gsessions.Session) {
	sess.Set(sessionNoChange, true)
}
