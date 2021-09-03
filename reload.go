package goWeb

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"goWeb/ext/weblib"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 待完善
func Exit() {
	rFile := "./temp.tmp"

	a, e := os.ReadFile(rFile)
	if e == nil {
		in(a)
	}

	f := func() {
		for {
			err := os.WriteFile(rFile, out(), 0666)
			if err != nil {
				log.Println("some out false     :", err)
			}
			time.Sleep(time.Minute / 1)
			// time.Sleep(time.Millisecond * 10)
		}
	}
	go f()
	//合建chan

	c := make(chan os.Signal, 100)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//阻塞直至有信号传入
	s := <-c
	fmt.Println("接收到退出信号", s, "正在退出")
	os.WriteFile("./temp.tmp", out(), 0666)
	os.Exit(0)
}

func out() []byte {
	a := make(map[string][]byte)
	a["ss"] = weblib.Sm.Out()

	return sOut(a)
}
func sOut(i interface{}) []byte {
	bu := new(bytes.Buffer)
	gob.NewEncoder(bu).Encode(i)
	return bu.Bytes()
}
func in(i []byte) {
	a := make(map[string][]byte)
	sIn(i, &a)
	weblib.Sm.In(a["ss"])
}
func sIn(i []byte, d interface{}) {
	bu := bytes.NewReader(i)
	gob.NewDecoder(bu).Decode(d)
}
