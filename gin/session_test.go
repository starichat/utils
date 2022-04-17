package gin

import (
	"math"
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	redisconn "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
)

func TestSession(t *testing.T) {
	eng := gin.Default()
	rstore, err := redis.NewStoreWithDB(10, "tcp", "127.0.0.1:6379", "requirepass", "2", []byte("secret"))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	rstore.Options(sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   1200,
		HttpOnly: true,
	})
	auto := &Store{pairs: securecookie.CodecsFromPairs([]byte("secret")), Store: rstore, AutoRefresh: true}
	eng.Group("").Use(sessions.Sessions("session", auto)).Any("auto", func(c *gin.Context) {
		sess := sessions.Default(c)
		key := uuid.New().String()
		if c.Query("save") != "" {
			sess.Set("key", key)
			sess.Save()
		} else {
			MarkSessionNoChange(sess)
			sess.Set("key", key)
			sess.Save()
		}
		c.String(http.StatusOK, key)
	})

	zero := &Store{pairs: securecookie.CodecsFromPairs([]byte("secret")), Store: rstore, MaxAgeZero: true}
	eng.Group("").Use(sessions.Sessions("session", zero)).Any("zero", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Delete("")
		sess.Save()
	})

	cons := &Store{pairs: securecookie.CodecsFromPairs([]byte("secret")), Store: rstore, ConstExp: true}
	eng.Group("").Use(sessions.Sessions("session", cons)).Any("const", func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Delete("")
		sess.Save()
	}).Any("const/nochange", func(c *gin.Context) {
		sess := sessions.Default(c)
		MarkSessionNoChange(sess) //标记状态(只进行续期)
		sess.Delete("")           //用于修改session的written状态，用于触发save方法
		sess.Save()               //save
	})

	go func() {
		eng.Run("0.0.0.0:18000")
	}()
	time.Sleep(time.Second * 2) //暂停2s，启动http服务
	testSessionAuto(t, rstore)
	testSessionZero(t, rstore)
	testSessionConst(t, rstore)
}

func testSessionAuto(t *testing.T, store redis.Store) {
	exp1, exp2, ttl1, ttl2, err := testTimeData("http://127.0.0.1:18000/auto", t, store)
	if err != nil {
		t.FailNow()
	}

	t.Log(exp1, exp2)
	t.Log(ttl1, ttl2)

	if math.Abs(float64(exp1.Sub(exp2))) < float64(time.Second) { //auto refresh session的时间会自动进行续期操作，sleep n秒，时间差值肯定要大于n-1秒
		t.FailNow()
	}

	if math.Abs(float64(ttl1-ttl2)) > float64(time.Millisecond*500) {
		t.FailNow()
	}
}

func testSessionZero(t *testing.T, store redis.Store) {
	exp1, exp2, ttl1, ttl2, err := testTimeData("http://127.0.0.1:18000/zero", t, store)
	if err != nil {
		t.FailNow()
	}

	t.Log(exp1, exp2)
	t.Log(ttl1, ttl2)

	if !exp1.IsZero() || !exp2.IsZero() {
		t.FailNow()
	}

	if math.Abs(float64(ttl1-ttl2)) > float64(time.Millisecond*500) {
		t.FailNow()
	}
}

func testSessionConst(t *testing.T, store redis.Store) {
	exp1, exp2, ttl1, ttl2, err := testTimeData("http://127.0.0.1:18000/const/nochange", t, store)
	if err != nil {
		t.FailNow()
	}

	t.Log(exp1, exp2)
	t.Log(ttl1, ttl2)
	if math.Abs(float64(exp1.Sub(exp2))) > float64(time.Second) {
		t.FailNow()
	}
	if math.Abs(float64(ttl1-ttl2-1000*2)) > float64(time.Second) {
		t.FailNow()
	}
}

func testTimeData(uri string, t *testing.T, store redis.Store) (exp1, exp2 time.Time, ttl1, ttl2 int64, err error) {
	//测试const
	client := &http.Client{}
	client.Jar, _ = cookiejar.New(nil)

	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	c := session(resp)
	cv := c.Value
	ttl1, err = getttl(store, cv)
	resp.Body.Close()
	if err != nil {
		return
	}
	time.Sleep(time.Second * 2) //延长时间

	resp, err = client.Get(uri)
	if err != nil {
		return
	}
	c1 := session(resp)
	cv = c1.Value
	ttl2, err = getttl(store, cv)
	resp.Body.Close()
	exp1 = c.Expires
	exp2 = c1.Expires
	return
}

func getttl(s redis.Store, v string) (int64, error) {
	id := ""
	securecookie.DecodeMulti("session", v, &id, securecookie.CodecsFromPairs([]byte("secret"))...)
	_, r := redis.GetRedisStore(s)
	conn := r.Pool.Get()
	defer conn.Close()
	return redisconn.Int64(conn.Do("PTTL", "session_"+id))
}

func session(resp *http.Response) *http.Cookie {
	cs := resp.Cookies()
	var c *http.Cookie
	for i := 0; i < len(cs); i++ {
		if cs[i].Name == "session" {
			c = cs[i]
		}
	}
	return c
}
