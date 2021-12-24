package ordercaller

import (
	"context"
	"sort"
)

//Caller  container for func
type Caller struct {
	Fn    func(context.Context)
	Score int
}

//NewCaller returns a new Caller
func NewCaller(fn func(context.Context), score int) *Caller {
	return &Caller{Fn: fn, Score: score}
}

//Callers Caller Collection
type Callers []*Caller

//NewCallers returns a new pointer of Callers
func NewCallers(cap int) *Callers {
	cs := make(Callers, 0, cap)
	return &cs
}

//Add to add Caller to Callers
func (cs *Callers) Add(c *Caller) {
	*cs = append(*cs, c)
}

//Value return Caller of specific index
func (cs Callers) Value(index int) *Caller {
	return cs[index]
}

//Sort implement Sort interface
func (cs Callers) Sort() {
	sort.Sort(cs)
}

//Len return length of Callers
func (cs Callers) Len() int {
	return len(cs)
}

//Less
func (cs Callers) Less(i, j int) bool {
	return cs[i].Score < cs[j].Score
}

//Swap
func (cs Callers) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}
