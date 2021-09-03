package controller

import (
	"github.com/cjie9759/goWeb/ext/session"
	"github.com/cjie9759/goWeb/ext/weblib"
)

type Index struct {
	BaseApp
	ss session.Session
}

func (t *Index) Init() {

	s := Sm.BeginSession(t.W, t.R)
	Sm.Update(t.W, t.R)
	t.ss = s
}
func (t *Index) SayHi() {
	weblib.NewWebBase(t.W, t.R).WebJson("hi")

}
