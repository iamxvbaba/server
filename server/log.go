package server

import (
	lg "log"
	"os"
)
var Log *lg.Logger
func init(){
	f, err := os.OpenFile("mw_server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if nil != err {
		panic(err)
	}
	Log = lg.New(f, "[mw]", lg.Ldate|lg.Ltime|lg.Lshortfile)
}