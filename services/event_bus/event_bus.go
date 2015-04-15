package event_bus

import (
	"log"
	"reflect"
	"sync"
)

type EventBus struct {
	handlers map[reflect.Type][]reflect.Value
	lock     sync.RWMutex
}

/**
 * Construct a new EventBus
 */
func New() *EventBus {
	return &EventBus{
		make(map[reflect.Type][]reflect.Value),
		sync.RWMutex{},
	}
}

/**
 * Register an event handler
 */
func (ebus *EventBus) Register(fn interface{}, forTypes ...interface{}) {
	v := reflect.ValueOf(fn)
	def := v.Type()

	if def.NumIn() != 1 {
		log.Panicf("EventBus Handler must have a single argument")
	}

	argument := def.In(0)

	for _, typ := range forTypes {
		t := reflect.TypeOf(typ)
		if !t.ConvertibleTo(argument) {
			log.Fatalf("EventBus Handler argument %v is not compatible with type %v", argument, t)
		}
		ebus.addHandler(t, v)
	}

	if len(forTypes) == 0 {
		ebus.addHandler(argument, v)
	}
}

/**
 * Publish an event to the EventBus
 */
func (ebus *EventBus) Publish(event interface{}) error {
	ebus.lock.RLock()
	defer ebus.lock.RUnlock()

	t := reflect.TypeOf(event)

	handlers, ok := ebus.handlers[t]
	if !ok {
		return nil
	}

	args := [...]reflect.Value{reflect.ValueOf(event)}
	for _, fn := range handlers {
		fn.Call(args[:])
	}
	return nil
}

func (ebus *EventBus) addHandler(fnType reflect.Type, fn reflect.Value) {
	ebus.lock.Lock()
	defer ebus.lock.Unlock()
	handlers, ok := ebus.handlers[fnType]
	if !ok {
		handlers = make([]reflect.Value, 0)
	}
	ebus.handlers[fnType] = append(handlers, fn)
}
