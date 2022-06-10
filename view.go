package goweb

import (
	"embed"

	"github.com/cjie9759/goWeb/Demo/controller"
	"github.com/cjie9759/goWeb/ext/weblib"
)

type View struct {
	controller.BaseApp
	// ss *ext.Session
	// Fs embed.FS
}

func (t *View) Init(fs *embed.FS) {
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
