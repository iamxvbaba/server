package main

import (
	"context"
	"fmt"
	"github.com/iamxvbaba/server/server"
	"github.com/iamxvbaba/server/upgrader"
	"net"
	"time"
)

type Server struct {
	ln net.Listener
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
	var err error
	s.ln, err = upg.Fds.Listen("tcp", "0.0.0.0:18541")
	if err != nil {
		return err
	}
	server.Log.Println("app Initialize")
	return nil
}

func (s *Server) Serve(ctx context.Context) {
	server.Log.Println("app Serve!!!!")
	s.normal()
	// s.tcpStart()
}

func (s *Server) Destroy() {
	server.Log.Println("app Destroy")
}

func (s *Server) Daemon() bool {
	return true
}
func (s *Server) normal() {
	server.Log.Println("normal server start!!!")
}
func (s *Server) tcpStart() {
	defer s.ln.Close()
	server.Log.Printf("listening on %s", s.ln.Addr())
	for {
		c, err := s.ln.Accept()
		if err != nil {
			server.Log.Printf("listening error:%v", err)
			continue
		}
		go func() {
			for {
				_, e := c.Write([]byte(fmt.Sprintf("app server response at %s\n",time.Now().Format("2006-01-02 15:04:04"))))
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