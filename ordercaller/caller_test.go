package ordercaller

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaller(t *testing.T) {
	c := Caller{}
	c.Score = 1

	c.Fn = func(c context.Context) {
		fmt.Println(123)
	}

	c.Fn(context.Background())
	assert.Equal(t, c.Score, 1)

}

func TestCallers(t *testing.T) {
	c1 := &Caller{}
	c1.Score = 2

	c2 := &Caller{}
	c2.Score = 1

	cs := Callers{}
	cs = append(cs, c1)
	cs = append(cs, c2)

	cs.Sort()

	assert.Equal(t, cs[0].Score, 1)
	assert.Equal(t, cs[1].Score, 2)

}
