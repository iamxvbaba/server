package main

import (
	"flag"
	"fmt"
	"github.com/iamxvbaba/server/upgrader"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		listenAddr = flag.String("listen", "localhost:8080", "`Address` to listen on")
		pidFile    = flag.String("pid-file", "./run_pid", "`Path` to pid file")
		//doUpgrade = flag.Bool("dup",false,"use upgrade")
	)

	flag.Parse()
	log.SetPrefix(fmt.Sprintf("%d ", os.Getpid()))

	upg, err := upgrader.New(upgrader.Options{
		PIDFile: *pidFile,
	})
	if err != nil {
		panic(err)
	}
	defer upg.Stop()

	// Do an upgrade on SIGHUP
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP)
		for range sig {
			fmt.Println("进行升级!!!!!!!")
			err := upg.Upgrade()
			if err != nil {
				log.Println("upgrade failed:", err)
			}
		}
	}()

	ln, err := upg.Fds.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalln("Can't listen:", err)
	}

	go func() {
		defer ln.Close()

		log.Printf("listening on %s", ln.Addr())

		for {
			c, err := ln.Accept()
			if err != nil {
				panic(err)
			}
			go func() {
				for {
					_, e := c.Write([]byte(fmt.Sprintf("Hello 2 at %s\n",time.Now().Format("2006-01-02 15:04:04"))))
					if e != nil {
						return
					}
					time.Sleep(5*time.Second)
				}
			}()
			go func() {
				data := make([]byte,128)
				_,e := c.Read(data)
				if e != nil {
					fmt.Println("断开：",e)
					return
				}
			}()
		}
	}()

	fmt.Println("ready!!!!!!!!!")
	if err := upg.Ready(); err != nil {
		panic(err)
	}
	<-upg.Exit()
}