package weblib

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/cjie9759/goWeb/ext/session"
)

var Sm *session.SessionManager

func init() {
	Sm = session.NewSessionManager()
}

type WebBase struct {
	r      *http.Request
	w      http.ResponseWriter
	isJson bool
	isErr  int
	code   int
	body   string
	head   map[string]string
}

func NewWebBase(w http.ResponseWriter, r *http.Request) *WebBase {
	a := &WebBase{
		r:      r,
		w:      w,
		isJson: true,
		isErr:  0,
		code:   200,
		body:   "",
		head:   make(map[string]string),
	}
	a.r.ParseForm()
	return a
}
func (B *WebBase) SetCode(c int) *WebBase {
	B.code = c
	return B
}
func (B *WebBase) SetBody(b string) *WebBase {
	B.body = b
	return B
}
func (B *WebBase) SetHead(k string, v string) *WebBase {
	B.head[k] = v
	return B
}
func (B *WebBase) IsErr() *WebBase {
	B.isErr = 1
	return B
}
func (B *WebBase) IsJson(b bool) *WebBase {
	B.isJson = b
	return B
}
func (B *WebBase) Send() {
	w := B.w
	for k, v := range B.head {
		B.w.Header().Set(k, v)
	}
	if B.isJson {
		w.Header().Set("content-type", "application/json; charset=utf-8")
		a := make(map[string]interface{})
		a["err"] = B.isErr
		a["msg"] = B.body
		b, _ := json.Marshal(a)
		w.WriteHeader(B.code)
		fmt.Fprintln(w, string(b))
		return
	}
	w.Header().Set("content-type", "text/html; charset=utf-8")
	w.WriteHeader(B.code)
	fmt.Fprintln(w, B.body)

}
func (B *WebBase) Get(name string) string {
	query := B.r.URL.Query()
	return query.Get(name)
}
func (B *WebBase) Post(name string) string {
	form := B.r.Form
	return form.Get(name)

}

func InArray(O *[]string, S string) bool {
	for _, v := range *O {
		if v == S {
			return true
		}
	}
	return false
}

func Md5(s1 string, s2 string) [16]byte {
	sign1 := md5.Sum([]byte(s1))
	sign2 := md5.Sum([]byte(s2))
	sign3 := md5.Sum(append(sign1[:], sign2[:]...))
	return sign3
}
func Pathinfo(r *http.Request, p1 string, p2 string) (string, string) {
	a := r.URL.Path
	b := strings.Split(a[3:], "/")
	switch len(b) {
	case 0, 1:
		return p1, p2
	case 2:
		return b[1], p2
	default:
		return b[1], b[2]
	}
}
func SetProxy(proxy string) {
	os.Setenv("HTTP_PROXY", proxy)
	os.Setenv("HTTPS_PROXY", proxy)
}
