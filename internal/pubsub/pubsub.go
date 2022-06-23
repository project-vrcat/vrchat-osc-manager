package pubsub

import (
	"context"
	"sync"
)

type (
	Hub[T any] struct {
		sync.Mutex
		cap  int
		subs map[*Sub[T]]struct{}
	}
	Sub[T any] struct {
		topic   string
		msg     chan *T
		quit    chan struct{}
		handler func(*T)
	}
)

func New[T any](cap int) *Hub[T] {
	return &Hub[T]{
		cap:  cap,
		subs: make(map[*Sub[T]]struct{}),
	}
}

func (h *Hub[T]) Publish(ctx context.Context, topic string, msg *T) {
	h.Lock()
	for s := range h.subs {
		if s.topic == topic {
			s.Publish(ctx, msg)
		}
	}
	h.Unlock()
}

func (h *Hub[T]) Sub(ctx context.Context, topic string, handler func(*T)) *Sub[T] {
	sub := &Sub[T]{
		topic:   topic,
		msg:     make(chan *T, h.cap),
		quit:    make(chan struct{}),
		handler: handler,
	}
	h.Lock()
	h.subs[sub] = struct{}{}
	h.Unlock()

	go func() {
		select {
		case <-sub.quit:
		case <-ctx.Done():
			h.Lock()
			delete(h.subs, sub)
			h.Unlock()
		}
	}()

	go sub.run(ctx)
	return sub
}

func (h *Hub[T]) UnSub(s *Sub[T]) {
	h.Lock()
	delete(h.subs, s)
	h.Unlock()
	close(s.quit)
}

func (s *Sub[T]) Publish(ctx context.Context, msg *T) {
	select {
	case <-ctx.Done():
		return
	case s.msg <- msg:
	default:
	}
}

func (s *Sub[T]) run(ctx context.Context) {
	for {
		select {
		case msg := <-s.msg:
			s.handler(msg)
		case <-s.quit:
			return
		case <-ctx.Done():
			return
		}
	}
}
