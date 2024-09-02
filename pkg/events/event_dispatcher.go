package events

import (
	"errors"
	"sync"
	"time"
)

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

type (
	Event interface {
		GetName() string
		GetDateTime() time.Time
		GetPayload() interface{}
		SetPayload(payload interface{})
	}

	EventHandler interface {
		Handle(event Event, wg *sync.WaitGroup)
	}

	EventDispatcher interface {
		Register(eventName string, handler EventHandler) error
		Dispatch(event Event) error
		Remove(eventName string, handler EventHandler) error
		Has(eventName string, handler EventHandler) bool
		Clear()
	}

	DefaultEventDispatcher struct {
		handlers map[string][]EventHandler
	}
)

func NewEventDispatcher() *DefaultEventDispatcher {
	return &DefaultEventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

func (ed *DefaultEventDispatcher) Dispatch(event Event) error {
	if handlers, ok := ed.handlers[event.GetName()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(event, wg)
		}
		wg.Wait()
	}
	return nil
}

func (ed *DefaultEventDispatcher) Register(eventName string, handler EventHandler) error {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *DefaultEventDispatcher) Has(eventName string, handler EventHandler) bool {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

func (ed *DefaultEventDispatcher) Remove(eventName string, handler EventHandler) error {
	if _, ok := ed.handlers[eventName]; ok {
		for i, h := range ed.handlers[eventName] {
			if h == handler {
				ed.handlers[eventName] = append(ed.handlers[eventName][:i], ed.handlers[eventName][i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (ed *DefaultEventDispatcher) Clear() {
	ed.handlers = make(map[string][]EventHandler)
}
