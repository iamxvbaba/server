package server

import (
	"context"
	"fmt"
	"github.com/iamxvbaba/server/upgrader"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type AppInstance interface {
	Name() string
	Version() string
	Initialize(context.Context) error
	Serve(ctx context.Context, upg *upgrader.Upgrader)
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
		upg         *upgrader.Upgrader
		ctx, cancel = context.WithCancel(context.Background())
	)
	rand.Seed(time.Now().UnixNano())
	if app == nil {
		panic("app instance is nil")
	}
	Log.SetPrefix(fmt.Sprintf("[app_%s_%d]", app.Name(), os.Getpid()))
	if err = app.Initialize(ctx); err != nil {
		panic(err)
	}
	Log.Printf("app:%s version:%s is running \n", app.Name(), app.Version())
	if upg, err = upgrader.New(upgrader.Options{
		PIDFile: fmt.Sprintf("%s_run_pid", app.Name()),
	}); err != nil {
		panic(err)
	}
	go app.Serve(ctx, upg)

	defer func() {
		Log.Printf("app:%s version:%s stop\n", app.Name(), app.Version())
		cancel()
		app.Destroy()
		upg.Stop()
	}()

	// Do an upgrade on SIGHUP
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGILL,syscall.SIGINT)
		switch x := <-sig; x {
		case syscall.SIGHUP:
			Log.Printf("app:%s 进行升级!!!!!!!", app.Name())
			err := upg.Upgrade()
			if err != nil {
				Log.Println("upgrade failed:", err)
			}
		default:
			Log.Printf("app:%s 退出 sigal:%v", app.Name(), x)
			upg.Stop()
		}
	}()
	if err := upg.Ready(); err != nil {
		panic(err)
	}
	<-upg.Exit()
}
