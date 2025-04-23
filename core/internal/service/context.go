package service

import (
	"context"
	"sync"
)

type Context struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewContext(ctx context.Context) Context {
	ctx, cancel := context.WithCancel(ctx)

	return Context{
		ctx:    ctx,
		cancel: cancel,
		wg:     sync.WaitGroup{},
	}
}

func (ctx *Context) Lock() {
	ctx.wg.Add(1)
}

func (ctx *Context) Unlock() {
	ctx.wg.Done()
}

func (ctx *Context) Wait() {
	ctx.wg.Wait()
}

func (ctx *Context) Cancel() {
	ctx.cancel()
}

func (ctx *Context) CancelAndWait() {
	ctx.cancel()
	ctx.wg.Wait()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.ctx.Done()
}

func (ctx *Context) Context() context.Context {
	return ctx.ctx
}
