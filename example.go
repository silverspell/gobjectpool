package gobjectpool

import (
	"fmt"
	"math/rand"
	"time"
)

func example() {
	p := new(Pool)
	rand.Seed(time.Now().Unix())
	p.OnCreateFunction = func() any {
		return rand.Intn(100)
	}

	p.OnDestroyFunction = func(a any) {
		a = 0
	}

	p.Init(&PoolOptions{
		MaxItems: 3,
	})

	val, _ := p.Borrow()
	fmt.Printf("Borrow: %+v\n", val.(int))
	val2, _ := p.Borrow()
	fmt.Printf("Borrow: %+v\n", val2.(int))
	val3, _ := p.Borrow()
	fmt.Printf("Borrow: %+v\n", val3.(int))
	p.Return(val3)
	val4, _ := p.Borrow()
	fmt.Printf("Borrow: %+v\n", val4.(int))
	val5, _ := p.Borrow()
	fmt.Printf("Borrow: %+v\n", val5.(int))
	p.Return(val)
	p.Return(val2)
	p.Return(val4)
	p.Return(val5)

	fmt.Printf("%+v\n", p.Items)
}
