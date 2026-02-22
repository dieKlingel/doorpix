package sip

import (
	"log/slog"
	"runtime"

	"github.com/dieklingel/go-pjproject/pjsua2"
)

type bridge struct {
	endpoint pjsua2.Endpoint
	counter  int

	channel chan cmd
	done    chan struct{}
}

var osThread *bridge

func init() {
	osThread = &bridge{
		channel: make(chan cmd),
		done:    make(chan struct{}),
	}
}

type cmd struct {
	f    func()
	done chan bool
}

func (ua *bridge) invoke(f func()) {
	// the endpoint is nil when invoke is called from the init function
	if ua.endpoint != nil && ua.endpoint.LibIsThreadRegistered() {
		f()
	}

	done := make(chan bool)
	ua.channel <- cmd{
		f:    f,
		done: done,
	}

	success := <-done
	close(done)

	if !success {
		slog.Error("the native os thread invocation did not finished successfully")
	}
}

func (ua *bridge) run() {
	runtime.LockOSThread()

	config := pjsua2.NewEpConfig()
	ua.endpoint = pjsua2.NewEndpoint()
	ua.endpoint.LibCreate()
	ua.endpoint.LibInit(config)
	ua.endpoint.LibStart()

	for {
		select {
		case cmd := <-ua.channel:
			cmd.f()
			cmd.done <- true
		case <-ua.done:
			return
		}
	}
}
