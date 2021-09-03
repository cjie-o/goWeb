package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var session_debug int

func init() {
	session_debug = 0
}

//-------------session_implements-----------------
//Session操作接口，不同存储方式的Sesion操作不同，实现也不同
type Session interface {
	Set(key, value interface{})
	Get(key interface{}) interface{}
	Remove(key interface{}) error
	GetId() string
	Update()
}

//--------------session_from---------------------------
//session存储方式接口，可以存储在内存，数据库或者文件
//分别实现该接口即可
//如存入数据库的CRUD操作
type Storage interface {
	//初始化一个session，id根据需要生成后传入
	InitSession(sid string, maxAge int64) (Session, error)
	//根据sid，获得当前session
	GetSession(sid string) Session
	//销毁session
	DestroySession(sid string) error
	//回收
	GCSession()
	In([]byte)
	Out() []byte
}

//--------------session_manager----------------------
//管理Session,实际操作cookie，Storage
//由于该结构体是整个应用级别的，写、修改都需要枷锁
type SessionManager struct {
	//session数据最终需要在客户端（浏览器）和服务器各存一份
	//客户端时，存放在cookie中
	cookieName string
	//存放方式，如内存，数据库，文件
	storage Storage
	//超时时间
	maxAge int64
	//由于session包含所有的请求
	//并行时，保证数据独立、一致、安全
	lock sync.Mutex
}

//实例化一个session管理器
func NewSessionManager() *SessionManager {
	sessionManager := &SessionManager{
		cookieName: "my-cookie",
		storage:    newFromMemory(), //默认以内存实现
		maxAge:     60 * 30,         //默认30分钟
	}
	go sessionManager.GC()

	return sessionManager
}

func (m *SessionManager) GetCookieN() string {
	return m.cookieName
}

//先判断当前请求的cookie中是否存在有效的session,存在返回，不存在创建
func (m *SessionManager) BeginSession(w http.ResponseWriter, r *http.Request) Session {
	//防止处理时，进入另外的请求
	m.lock.Lock()
	defer m.lock.Unlock()

	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" { //如果当前请求没有改cookie名字对应的cookie
		if session_debug == 1 {
			log.Println("session debug", "-----------> current session not exists")
		}
		//创建一个
		sid := m.randomId()
		//根据保存session方式，如内存，数据库中创建
		session, _ := m.storage.InitSession(sid, m.maxAge) //该方法有自己的锁，多处调用到

		maxAge := m.maxAge

		//用session的ID于cookie关联
		//cookie名字和失效时间由session管理器维护
		cookie := http.Cookie{
			Name: m.cookieName,
			//这里是并发不安全的，但是这个方法已上锁
			Value:    url.QueryEscape(sid), //转义特殊符号@#￥%+*-等
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(maxAge),
			Expires:  time.Now().Add(time.Duration(maxAge)),
		}
		http.SetCookie(w, &cookie) //设置到响应中
		return session
	} else { //如果存在

		sid, _ := url.QueryUnescape(cookie.Value) //反转义特殊符号
		session := m.storage.GetSession(sid)      //从保存session介质中获取
		// session := m.storage.
		if session_debug == 1 {
			log.Println("session debug", "session --------->", session)
		}
		if session == nil {
			if session_debug == 1 {
				log.Println("session debug", "-----------> current session is nil")
			}
			//创建一个
			sid := m.randomId()
			//根据保存session方式，如内存，数据库中创建
			newSession, _ := m.storage.InitSession(sid, m.maxAge) //该方法有自己的锁，多处调用到

			maxAge := m.maxAge

			//用session的ID于cookie关联
			//cookie名字和失效时间由session管理器维护
			newCookie := http.Cookie{
				Name: m.cookieName,
				//这里是并发不安全的，但是这个方法已上锁
				Value:    url.QueryEscape(sid), //转义特殊符号@#￥%+*-等
				Path:     "/",
				HttpOnly: true,
				MaxAge:   int(maxAge),
				Expires:  time.Now().Add(time.Duration(maxAge)),
			}
			http.SetCookie(w, &newCookie) //设置到响应中
			return newSession
		}
		// else {
		// 	m.Update(w, r)
		// }
		if session_debug == 1 {
			log.Println("session debug", "-----------> current session exists")
		}
		return session
	}

}

//更新超时
func (m *SessionManager) Update(w http.ResponseWriter, r *http.Request) {
	m.lock.Lock()
	defer m.lock.Unlock()

	cookie, err := r.Cookie(m.cookieName)
	if err != nil {
		return
	}
	sid, _ := url.QueryUnescape(cookie.Value)

	// fmt.Println(m)
	fmt.Println(sid)
	s := m.storage.GetSession(sid)
	if s == nil {
		return
	}
	s.Update()
	cookie.Path = "/"
	cookie.MaxAge = int(m.maxAge)
	http.SetCookie(w, cookie)
}

//通过ID获取session
func (m *SessionManager) GetSessionById(sid string) Session {
	session := m.storage.GetSession(sid)
	return session
}

//手动销毁session，同时删除cookie
func (m *SessionManager) Destroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		m.lock.Lock()
		defer m.lock.Unlock()

		sid, _ := url.QueryUnescape(cookie.Value)
		m.storage.DestroySession(sid)

		cookie2 := http.Cookie{
			MaxAge:  0,
			Name:    m.cookieName,
			Value:   "",
			Path:    "/",
			Expires: time.Now().Add(time.Duration(0)),
		}

		http.SetCookie(w, &cookie2)
	}
}

func (m *SessionManager) CookieIsExists(r *http.Request) bool {
	_, err := r.Cookie(m.cookieName)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (m *SessionManager) GC() {
	t1 := time.NewTicker(time.Duration(m.maxAge / 10))
	f := func() {
		m.lock.Lock()
		defer m.lock.Unlock()

		m.storage.GCSession()
	}
	for {
		<-t1.C
		f()
	}
}

func (m *SessionManager) SetMaxAge(t int64) {
	m.maxAge = t
}

func (m *SessionManager) SetSessionFrom(storage Storage) {
	m.storage = storage
}

func (m *SessionManager) randomId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	//加密
	return base64.URLEncoding.EncodeToString(b)
}
func (m *SessionManager) Out() []byte {
	return m.storage.Out()
}
func (m *SessionManager) In(i []byte) {
	m.storage.In(i)
}
