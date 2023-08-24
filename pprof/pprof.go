package pprof

import (
	"net/http"
	_ "net/http/pprof"
)

func InitPprof() {
	go func() {
		http.ListenAndServe("0.0.0.0:19999", nil)
	}()
}
