package subpub

import (
	"context"
	"sync"
	"sync/atomic"
)

type Subscriber struct {
	ch     chan interface{}
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	closed atomic.Bool
}
