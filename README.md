# golang程序启动管理


##### 参考项目: https://github.com/cloudflare/tableflip.git

##### example

```
package main

import (
	"context"
	"fmt"
	"github.com/iamxvbaba/server/server"
	"github.com/iamxvbaba/server/upgrader"
	"time"
)

func main() {
	server.Run(New())
}

type Server struct {

}

func New() *Server{
	return &Server{}
}

func (s *Server) Name() string {
	return "test_server"
}

func (s *Server) Version() string {
	return "v0.0.2"
}

func (s *Server) Initialize(ctx context.Context, upg *upgrader.Upgrader) error {
	server.Log.Println("app Initialize")
	return nil
}

func (s *Server) Serve(ctx context.Context, upg *upgrader.Upgrader) {
	server.Log.Println("app Serve!!!!")
	ln, err := upg.Fds.Listen("tcp", "0.0.0.0:18541")
	if err != nil {
		server.Log.Fatalln("Can't listen:", err)
	}
	defer ln.Close()
	server.Log.Printf("listening on %s", ln.Addr())
	for {
		c, err := ln.Accept()
		if err != nil {
			server.Log.Printf("listening error:%v", err)
			continue
		}
		go func() {
			for {
				_, e := c.Write([]byte(fmt.Sprintf("response at %s\n",time.Now().Format("2006-01-02 15:04:04"))))
				if e != nil {
					return
				}
				time.Sleep(3*time.Second)
			}
		}()
		go func() {
			data := make([]byte,128)
			_,e := c.Read(data)
			if e != nil {
				server.Log.Println("断开：",e)
				return
			}
			server.Log.Printf("client data:%s",data)
		}()
	}
}

func (s *Server) Destroy() {
	server.Log.Println("app Destroy")
}

func (s *Server) Daemon() bool {
	return true
}


```

