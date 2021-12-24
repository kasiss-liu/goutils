package ordercaller

import (
	"context"
	"time"
)

//ChainContext to control Caller Chain
type ChainContext struct {
	wait chan struct{}
	ctx  context.Context
}

//NewChainContext return a new pointer of ChainContext
func NewChainContext(ctx ...context.Context) *ChainContext {
	c := &ChainContext{}
	c.wait = make(chan struct{}, 1)
	c.wait <- struct{}{}
	if len(ctx) > 0 && ctx[0] != nil {
		c.ctx = ctx[0]
	} else {
		c.ctx = context.TODO()
	}
	return c
}

//Deadline implement interface Context
func (cc ChainContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

//Done implement interface Context
func (cc ChainContext) Done() <-chan struct{} {
	return cc.wait
}

//Value implement interface Context
func (cc ChainContext) Value(key interface{}) interface{} {
	return cc.ctx.Value(key)
}

//Err implement interface Context
func (cc ChainContext) Err() error {
	return nil
}

//Wait  implement interface Context
func (cc ChainContext) Wait() {
	cc.wait <- struct{}{}
}

//Chain is a machine for callers
type Chain struct {
	Callers *Callers
	Ctx     *ChainContext
}

//Append Caller to chain
func (ch *Chain) Append(c *Caller) {
	ch.Callers.Add(c)
}

//Run callers in chain
func (ch *Chain) Run() {
	ch.Callers.Sort()
	l := ch.Callers.Len()

	for i := 0; i < l; i++ {
		go ch.Callers.Value(i).Fn(ch.Ctx)
		ch.Ctx.Wait()
	}
}

//NewFnChain
func NewFnChain(cap ...int) *Chain {
	c := 5
	if len(cap) > 0 {
		c = cap[0]
	}
	ch := &Chain{}
	ch.Callers = NewCallers(c)
	ch.Ctx = NewChainContext()
	return ch
}
