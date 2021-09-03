package goWeb

import (
	"fmt"
	"goWeb/ext/weblib"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"runtime/debug"
	"time"
)

type application struct {
	routes  map[string]reflect.Type
	isDebug bool
	m       func(h http.Handler) http.Handler
}

func NewApp() *application {
	return &application{
		routes:  make(map[string]reflect.Type),
		isDebug: false,
		m:       func(h http.Handler) http.Handler { return h },
	}
}

func (p *application) Debug() *application {
	p.isDebug = true
	return p
}
func (p *application) SetMiddle(f func(h http.Handler) http.Handler) *application {
	p.m = f
	return p
}

func (p *application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wB := weblib.NewWebBase(w, r)
	defer func() {
		if e := recover(); e != nil {
			wB.Web500()
			log.Println(e)
			if p.isDebug {
				log.Println(string(debug.Stack()))
			}
		}
	}()

	url1 := r.URL.String()
	is_api, _ := regexp.MatchString("/api.*", url1)
	if !is_api {
		a := &View{}
		a.W = w
		a.R = r
		a.Init()
		return
	}

	p1, p2 := weblib.Pathinfo(r, "", "")

	if p1 == "" || p2 == "" {
		wB.Web404()
		return
	}
	route, ok := p.routes[p1]
	if !ok {
		wB.Web404()
		return
	}

	call := func() {
		ele := reflect.New(route)
		ele.Elem().FieldByName("R").Set(reflect.ValueOf(r))
		ele.Elem().FieldByName("W").Set(reflect.ValueOf(w))
		if ele.MethodByName(p2).Kind().String() == "invalid" {
			wB.Web404()
			return
		}
		ele.MethodByName("Init").Call(nil)
		ele.MethodByName(p2).Call(nil)
	}
	call()
}

func (p *application) Get(c interface{}) *application {
	ele := reflect.TypeOf(c).Elem()
	p.routes[ele.Name()] = ele
	return p
}

func (p *application) Run(addr string) error {
	p.printRoutes()
	fmt.Printf("listen on %s\n", addr)
	return http.ListenAndServe(addr, p.m(p))
}

func (p *application) printRoutes() {
	for k, v := range p.routes {
		n := reflect.New(v).Type()
		fmt.Println(k, n.String())
		for i := 0; i < n.NumMethod(); i++ {
			fmt.Println("	func", n.Method(i).Name)
		}
	}
}

func MWLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		tn := time.Now().UnixNano()
		h.ServeHTTP(rw, r)
		log.Println("webLog", r.Host, r.Method, r.Proto, r.RemoteAddr,
			r.RequestURI, (time.Now().UnixNano()-tn)/1e6)
	})
}
