package main

import (
	"context"
	"fmt"
	"github.com/iamxvbaba/server/upgrader"
	"log"
	"time"
)

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

func (s *Server) Initialize(ctx context.Context) error {
	log.Println("app Initialize")
	return nil
}

func (s *Server) Serve(ctx context.Context, upg *upgrader.Upgrader) {
	log.Println("app Serve!!!!")
	ln, err := upg.Fds.Listen("tcp", "0.0.0.0:18541")
	if err != nil {
		log.Fatalln("Can't listen:", err)
	}
	defer ln.Close()
	log.Printf("listening on %s", ln.Addr())
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Printf("listening error:%v", err)
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
				log.Println("断开：",e)
				return
			}
			log.Printf("client data:%s",data)
		}()
	}
}

func (s *Server) Destroy() {
	log.Println("app Destroy")
}

func (s *Server) Daemon() bool {
	return true
}

