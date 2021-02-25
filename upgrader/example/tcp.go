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
		pidFile    = flag.String("pid-file", "", "`Path` to pid file")
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
				fmt.Println("conn close:",err)
				return
			}
			go func() {
				for {
					c.SetDeadline(time.Now().Add(time.Second))
					c.Write([]byte(fmt.Sprintf("Hello at %s\n",time.Now().Format("2006-01-02 15:04:04"))))
					time.Sleep(5*time.Second)
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