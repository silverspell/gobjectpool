package gobjectpool

import (
	"fmt"
	"sync"
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
}

func (p *Pool) Return(item any) {
	p.Locker.Lock()
	defer p.Locker.Unlock()
	if p.HasExceededMax() {
		fmt.Println("Found a rogue one, cleaning")
		p.destroy(item)
		return
	}
	p.Items = append(p.Items, item)
}

func (p *Pool) destroy(item any) {
	p.OnDestroyFunction(item)
}

func (p *Pool) Borrow() (any, error) {

	p.Locker.Lock()
	defer p.Locker.Unlock()
	if p.IsEmpty() { // in case pool is exhausted, we will try to add new item with OnCreate, on return we will check and destroy the mess
		p.Items = append(p.Items, p.OnCreateFunction())
	}
	item := p.Items[0]
	p.Items = p.Items[1:]
	return item, nil
}

func (p *Pool) IsEmpty() bool {
	return len(p.Items) == 0 //we can not lock here, caller has already a lock
}

func (p *Pool) HasExceededMax() bool {
	return uint32(len(p.Items)) > p.MaxItems // We also can not lock. Caller has already locked
}
