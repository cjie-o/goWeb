package controller

import (
	"goWeb/ext/session"
	"goWeb/ext/weblib"
	"net/http"
)

var Sm *session.SessionManager

func init() {
	Sm = weblib.Sm
}

type BaseApp struct {
	W http.ResponseWriter
	R *http.Request
}
