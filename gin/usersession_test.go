package gin

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	redisgo "github.com/gomodule/redigo/redis"
	"github.com/gorilla/securecookie"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"
	"time"
)


func TestUserSessionStore_DestroySessionByOpenID(t *testing.T) {
	var userStore *UserSessionStore
	type User struct {
		ID   string
		Name string
	}
	gob.Register(&User{})
	eng := gin.Default()
	rstore, err := redis.NewStoreWithDB(10, "tcp", "redis-master.sndu.cn:6379", "deepin!@#", "23", []byte("secret"))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	rstore.Options(sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   10,
		HttpOnly: false,
	})
	auto := &Store{pairs: securecookie.CodecsFromPairs([]byte("secret")), Store: rstore, AutoRefresh: true}
	userStore = &UserSessionStore{
		Store:         auto,
		UserKeyPrefix: "usersession",
	}

	eng.Group("").Use(sessions.Sessions("session", userStore)).Any("testlogin", func(c *gin.Context) {
		sess := sessions.Default(c)
		if c.Query("login") != "" {
			var u User
			u.ID = c.Query("login")
			u.Name = "loginName_" + c.Query("login")
			//保存用户，模拟登录操作
			sess.Set("user", &u)
			sess.Set(userStore.UserKeyPrefix, strings.TrimSpace(c.Query("login")))
			if err := sess.Save(); err != nil {
				t.Log(err)
				t.FailNow()
			}
			t.Log(sess.Get("user"))
			t.Log(sess.Get(userStore.UserKeyPrefix))
			c.String(http.StatusOK, "login "+c.Query("login")+" success")

		} else if c.Query("logout") != "" {
			//模拟退出登录
			sess.Options(sessions.Options{MaxAge: -1})
			sess.Clear()
			sess.Save()
			c.String(http.StatusOK, "logout "+c.Query("logout")+" success")

		} else if c.Query("watch") != "" {

			//查看用户信息
			userInfo := sess.Get("user")
			fmt.Println(userInfo)
			sess.Flashes()
			sess.Save()
			if u, ok := userInfo.(*User); ok && c.Query("watch") == u.ID {
				//未授权
				c.String(http.StatusOK, "watch "+c.Query("watch")+" success")
				return
			} else {
				c.Status(http.StatusUnauthorized)
			}

		} else if c.Query("unregister") != "" {
			t.Log(userStore.DestroyUserSessions(c.Query("unregister")))
			c.String(http.StatusOK, "unregister "+c.Query("unregister")+" success")

		} else if c.Query("loginnouser") != "" {
			//仅登录，不关联用户信息
			var u User
			u.ID = c.Query("login")
			u.Name = "loginName_" + c.Query("login")
			//保存用户，模拟登录操作
			sess.Set("user", u)
			if err := sess.Save(); err != nil {
				t.Log(err)
				t.FailNow()
			}
			t.Log(sess.Get("user"))
			c.String(http.StatusOK, "loginnouser "+c.Query("loginnouser")+" success")
		}
	})

	eng.Run("0.0.0.0:18000")

}

//登录授权测试
func TestLogin(t *testing.T) {
	resp, err := http.DefaultClient.Get("http://0.0.0.0:18000/testlogin?login=ae4e6b70c2f24771bf5d33138a258226")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	//获取对应的cookie，模拟对应的cookie进行请求，获取对应的登录信息
	cookie := session(resp)
	t.Log(cookie)
	//time.Sleep(9 * time.Second)
	code := testUserWatch(t, cookie, "ae4e6b70c2f24771bf5d33138a258226")
	if code != http.StatusOK {
		t.Log(code)
		t.FailNow()
	}

	testUserLogout(t, cookie,"ae4e6b70c2f24771bf5d33138a258226")

	// 去除cookie
	//req1, _ := http.NewRequest("GET", "http://0.0.0.0:18000/testlogin?watch=ae4e6b70c2f24771bf5d33138a258226", nil)
	//r1, err := http.DefaultClient.Do(req1)
	//if err != nil {
	//	t.Log(err)
	//	t.FailNow()
	//}
	//if r1.StatusCode != http.StatusUnauthorized {
	//	t.Log(r)
	//	t.FailNow()
	//}
	//t.Log("step2", r.Status)

}

func TestLoginNoUser(t *testing.T) {
	resp, err := http.DefaultClient.Get("http://0.0.0.0:18000/testlogin?login=ae4e6b70c2f24771bf5d33138a258226")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	//获取对应的cookie，模拟对应的cookie进行请求，获取对应的登录信息
	cookie := session(resp)
	t.Log(cookie)
	//模拟登录
	req, _ := http.NewRequest("GET", "http://0.0.0.0:18000/testlogin?login=ae4e6b70c2f24771bf5d33138a258226", nil)

	req.AddCookie(cookie)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if r.StatusCode != http.StatusOK {
		t.Log(r)
		t.FailNow()
	}
	t.Log("step1", r.Status)

	// 去除cookie
	req1, _ := http.NewRequest("GET", "http://0.0.0.0:18000/testlogin?watch=ae4e6b70c2f24771bf5d33138a258226", nil)
	r1, err := http.DefaultClient.Do(req1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if r1.StatusCode != http.StatusUnauthorized {
		t.Log(r)
		t.FailNow()
	}
	t.Log("step2", r.Status)
}

func TestMulSessionLogin(t *testing.T) {
	//同一个用户登录5台机器
	userID := "e2cdc22c-ab31-475e-a621-40dc6b3ca999"
	UserNumber := 10
	cookies := make([]*http.Cookie, UserNumber)
	for i := 0; i < UserNumber; i++ {
		cookies[i] = testUserLogin(t, userID)
	}
	//查看该用户当前有效session有多少个
	if testGetSessionStatus(t, userID) != UserNumber {
		t.Log("step1 ", testGetSessionStatus(t, userID))
		t.FailNow()
	}
	//第二阶段，退出登录两个个终端，查看session有多少个
	testUserLogout(t, cookies[1], userID)
	testUserLogout(t, cookies[4], userID)
	if testGetSessionStatus(t, userID) != UserNumber-2 {
		t.Log("step2 ", testGetSessionStatus(t, userID))
		t.FailNow()
	}

	time.Sleep(5 * time.Second)
	if code := testUserWatch(t, cookies[2], userID); code != http.StatusOK {
		t.Log("step2 testUserWatch", userID)
		t.FailNow()
	}
	if code := testUserWatch(t, cookies[3], userID); code != http.StatusOK {
		t.Log("step2 testUserWatch", userID)
		t.FailNow()
	}
	if code := testUserWatch(t, cookies[5], userID); code != http.StatusOK {
		t.Log("step2 testUserWatch", userID)
		t.FailNow()
	}
	if code := testUserWatch(t, cookies[8], userID); code != http.StatusOK {
		t.Log("step2 testUserWatch", userID)
		t.FailNow()
	}
	if code := testUserWatch(t, cookies[9], userID); code != http.StatusOK {
		t.Log("step2 testUserWatch", userID)
		t.FailNow()
	}

	//第三阶段，等待其中两个终端过期，再查看当前session有多少个
	time.Sleep(6 * time.Second)
	if testGetSessionStatus(t, userID) != 5 {
		t.Log("step3 ", testGetSessionStatus(t, userID))
		t.FailNow()
	}

	//第四阶段，登录另外一个用户
	secondUser := "9d531e71-1ee8-4274-bcbe-ece446622999"
	secondCookie := testUserLogin(t, secondUser)

	//第5阶段，注销前一个用户，并查看这俩个用户的状态是否符合预期
	testUserUnregister(t, userID)
	if testGetSessionStatus(t, userID) != 0 {
		t.Log("step5 ", testGetSessionStatus(t, userID))
		t.FailNow()
	}
	if code := testUserWatch(t, secondCookie, secondUser); code != 200 {
		t.Log("step6 ", testGetSessionStatus(t, secondUser))
		t.FailNow()
	}

}

func testUserLogin(t *testing.T, userID string) *http.Cookie {
	req, _ := http.NewRequest("GET", "http://0.0.0.0:18000/testlogin?login="+userID, nil)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log("login UserID" + userID + "success")
	return session(r)
}

func testUserLogout(t *testing.T, cookie *http.Cookie, userID string) {
	//req, _ := http.NewRequest("GET", "http://0.0.0.0:18000/testlogin?logout="+userID, nil)
	//req.AddCookie(cookie)

	url, _ := url.Parse("http://0.0.0.0:18000/")
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(url, []*http.Cookie{cookie}) //这里的cookies是[]*http.Cookie
	cli := http.Client{Transport: nil, Jar: jar}

	r, err := cli.Get("http://0.0.0.0:18000/testlogin?logout=" + userID)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if r.StatusCode != http.StatusOK {
		t.Log(r)
		t.FailNow()
	}

	t.Log("logout UserID" + userID + "success")
}

func testUserWatch(t *testing.T, cookie *http.Cookie, userID string) int {
	//req, _ := http.NewRequest("GET", "http://0.0.0.0:18000/testlogin?watch="+userID, nil)
	//req.AddCookie(cookie)
	//r, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	t.Log(err)
	//	t.FailNow()
	//}
	url, _ := url.Parse("http://0.0.0.0:18000/")
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(url, []*http.Cookie{cookie}) //这里的cookies是[]*http.Cookie
	cli := http.Client{Transport: nil, Jar: jar}

	r, err := cli.Get("http://0.0.0.0:18000/testlogin?watch=" + userID)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log("watch UserID" + userID + "success")
	t.Log(r)
	return r.StatusCode
}

func testUserUnregister(t *testing.T, userID string) {
	req, _ := http.NewRequest("GET", "http://0.0.0.0:18000/testlogin?unregister="+userID, nil)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if r.StatusCode != http.StatusOK {
		t.Log(r)
		t.FailNow()
	}

	t.Log("unregister UserID" + userID + "success")
}

// 获取当前用户有效sessionid的数量
func testGetSessionStatus(t *testing.T, userID string) int {
	var count int
	rstore, err := redis.NewStoreWithDB(10, "tcp", "redis-master.sndu.cn:6379", "deepin!@#", "23", []byte("secret"))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	err, rs := redis.GetRedisStore(rstore)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	conn := rs.Pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		t.Log(err)
		t.FailNow()
	}
	//获取所有的sessionID
	rpy, err := redisgo.Values(conn.Do("HKEYS", fmt.Sprintf("%s_%s", "usersession_", userID)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	fields := make([]string, len(rpy))
	for k, v := range rpy {
		if b, ok := v.([]byte); ok {
			fields[k] = string(b)
		}
	}
	//判断sessionid是否过期了
	for _, v := range fields {
		t.Log(v)
		ttl, err := redisgo.Int64(conn.Do("PTTL", "session_"+v))
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		if ttl > 0 {
			count++
		}
	}
	return count
}
