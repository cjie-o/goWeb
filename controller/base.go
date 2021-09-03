package controller

import (
	"net/http"

	"github.com/cjie9759/goWeb/ext/session"
	"github.com/cjie9759/goWeb/ext/weblib"
)

var Sm *session.SessionManager

func init() {
	Sm = weblib.Sm
}

type BaseApp struct {
	W http.ResponseWriter
	R *http.Request
}
