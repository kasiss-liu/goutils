package ordercaller

import (
	"context"
	"testing"
)

func TestChainRun(t *testing.T) {
	ch := NewFnChain(5)
	c5 := NewCaller(func(c context.Context) { println(5); <-c.Done() }, 5)
	ch.Append(c5)
	c4 := NewCaller(func(c context.Context) { println(4); <-c.Done() }, 4)
	ch.Append(c4)
	c3 := NewCaller(func(c context.Context) { println(3); <-c.Done() }, 3)
	ch.Append(c3)
	c2 := NewCaller(func(c context.Context) { println(2); <-c.Done() }, 2)
	ch.Append(c2)
	c1 := NewCaller(func(c context.Context) { println(1); <-c.Done() }, 1)
	ch.Append(c1)

	ch.Run()
}
