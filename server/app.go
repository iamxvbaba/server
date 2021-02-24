package server

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

type AppInstance interface {
	Name() string
	Version() string
	Initialize(context.Context) error
	Serve(context.Context)
	Destroy()
	Daemon() bool
}

func Run(app AppInstance) {
	if innerProcess {
		start(app)
	} else {
		if app.Daemon() {
			daemon()
		} else {
			start(app)
		}
	}

}
func start(app AppInstance) {
	var (
		err         error
		ctx, cancel = context.WithCancel(context.Background())
	)
	rand.Seed(time.Now().UnixNano())
	if app == nil {
		panic("app instance is nil")
	}
	if err = app.Initialize(ctx); err != nil {
		panic(err)
	}
	fmt.Printf("app:%s version:%s is running \n", app.Name(), app.Version())
	go app.Serve(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, sigTerm)

	s := <-ch
	fmt.Printf("app:%s version:%s exit by signal:%v \n", app.Name(), app.Version(), s)

	cancel()
	app.Destroy()

	fmt.Printf("app:%s version:%s exit \n", app.Name(), app.Version())
}
