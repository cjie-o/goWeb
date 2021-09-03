package session

import (
	"bytes"
	"encoding/gob"
	"log"
	"sync"
	"time"
)

//session实现
type SessionFromMemory struct {
	Sid              string
	lock             sync.Mutex
	LastAccessedTime time.Time
	MaxAge           int64
	Data             map[interface{}]interface{}
}

func newSessionFromMemory() *SessionFromMemory {
	return &SessionFromMemory{
		Data:   make(map[interface{}]interface{}),
		MaxAge: 60 * 30,
	}
}
func (si *SessionFromMemory) Set(key, value interface{}) {
	si.lock.Lock()
	defer si.lock.Unlock()
	si.Data[key] = value
}
func (si *SessionFromMemory) Get(key interface{}) interface{} {
	if value := si.Data[key]; value != nil {
		return value
	}
	return nil
}
func (si *SessionFromMemory) Remove(key interface{}) error {
	if value := si.Data[key]; value != nil {
		delete(si.Data, key)
	}
	return nil
}
func (si *SessionFromMemory) GetId() string {
	return si.Sid
}
func (si *SessionFromMemory) Update() {
	si.LastAccessedTime = time.Now()
}

//session来自内存
type FromMemory struct {
	lock     sync.Mutex
	Sessions map[string]Session
}

func newFromMemory() *FromMemory {
	return &FromMemory{
		Sessions: make(map[string]Session),
	}
}

func (fm *FromMemory) InitSession(Sid string, MaxAge int64) (Session, error) {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	newSession := newSessionFromMemory()
	newSession.Sid = Sid
	if MaxAge != 0 {
		newSession.MaxAge = MaxAge
	}
	newSession.LastAccessedTime = time.Now()

	fm.Sessions[Sid] = newSession
	return newSession, nil
}
func (fm *FromMemory) GetSession(Sid string) Session {
	return fm.Sessions[Sid]
}
func (fm *FromMemory) DestroySession(Sid string) error {
	if _, ok := fm.Sessions[Sid]; ok {
		delete(fm.Sessions, Sid)
		return nil
	}
	return nil
}
func (fm *FromMemory) GCSession() {

	Sessions := fm.Sessions

	//if session_debug == 1 { log.Println("session debug","gc session")}

	if len(Sessions) < 1 {
		return
	}

	//if session_debug == 1 { log.Println("session debug","current active Sessions ", Sessions)}

	for k, v := range Sessions {
		t := (v.(*SessionFromMemory).LastAccessedTime.Unix()) + (v.(*SessionFromMemory).MaxAge)

		if t < time.Now().Unix() {
			if session_debug == 1 {
				log.Println("session debug", "timeout-------->", v)
			}
			delete(fm.Sessions, k)
		}
	}

}
func (m *FromMemory) Out() []byte {
	m.lock.Lock()
	defer m.lock.Unlock()

	bu := new(bytes.Buffer)
	a := make(map[string]*SessionFromMemory)
	for k, v := range m.Sessions {
		a[k] = v.(*SessionFromMemory)
	}
	gob.NewEncoder(bu).Encode(a)
	return bu.Bytes()
}
func (m *FromMemory) In(i []byte) {
	m.lock.Lock()
	defer m.lock.Unlock()

	bu := bytes.NewReader(i)
	a := make(map[string]*SessionFromMemory)
	gob.NewDecoder(bu).Decode(&a)
	for k, v := range a {
		m.Sessions[k] = v
	}
}
