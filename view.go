package goWeb

import (
	"embed"

	"github.com/cjie9759/goWeb/controller"
	"github.com/cjie9759/goWeb/ext/weblib"
)

//go:embed public/*
var fs embed.FS

type View struct {
	controller.BaseApp
	// ss *ext.Session
	// Fs embed.FS
}

func (t *View) Init() {
	a := t.R.URL.String()
	b, err := fs.ReadFile("public" + a)
	if err != nil {
		if a == "/" {
			weblib.NewWebBase(t.W, t.R).WebLocation("/index.html")
			return
		}
	}
	t.W.Write(b)
}
