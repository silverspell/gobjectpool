package gobjectpool

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type CreateFunction func() any
type DestroyFunction func(any)
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) UnLock() {}

type PoolOptions struct {
	MaxItems uint32
}

type Pool struct {
	noCopy            noCopy
	MaxItems          uint32
	CurrentItemCount  uint32
	Locker            sync.Mutex
	Items             []any
	OnCreateFunction  CreateFunction
	OnDestroyFunction DestroyFunction
}

func (p *Pool) Init(options *PoolOptions) {
	p.MaxItems = options.MaxItems
	p.Items = make([]any, p.MaxItems)
	for i := 0; i < int(p.MaxItems); i++ {
		p.Items[i] = p.OnCreateFunction()
	}
	atomic.StoreUint32(&p.CurrentItemCount, p.MaxItems)
}

func (p *Pool) Return(item any) {
	if p.HasExceededMax() {
		fmt.Println("Found a rogue one, cleaning")
		p.OnDestroyFunction(item)
		atomic.AddUint32(&p.CurrentItemCount, ^uint32(0))
		return
	}
	p.Locker.Lock()
	p.Items = append(p.Items, item)
	p.Locker.Unlock()
	atomic.AddUint32(&p.CurrentItemCount, 1)
}

func (p *Pool) Borrow() (any, error) {

	p.Locker.Lock()
	defer p.Locker.Unlock()
	if p.IsEmpty() { // in case pool is exhausted, we will try to add new item with OnCreate, on return we will check and destroy the mess
		p.Items = append(p.Items, p.OnCreateFunction())
		atomic.StoreUint32(&p.CurrentItemCount, 1)
	}
	item := p.Items[0]
	p.Items = p.Items[1:]
	atomic.AddUint32(&p.CurrentItemCount, ^uint32(0))
	return item, nil
}

func (p *Pool) IsEmpty() bool {
	return p.CurrentItemCount == 0
}

func (p *Pool) HasExceededMax() bool {
	return p.CurrentItemCount > p.MaxItems
}
