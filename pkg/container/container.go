// Package container provides an IoC container for Go projects.
// It provides simple, fluent and easy-to-use interface to make dependency injection in GoLang easier.
package container

import (
	"reflect"
	"sync"
)

// binding keeps a binding resolver and instance (for singleton bindings).
type binding struct {
	resolver  interface{} // resolver function
	singleton bool

	mu       sync.Mutex
	instance interface{} // instance stored for singleton bindings
}

// Container is a map of reflect.Type to binding
type Container struct {
	parent *Container

	mu       sync.RWMutex
	bindings map[reflect.Type]map[string]*binding
}

// NewContainer returns a new instance of Container
func NewContainer() Container {
	return Container{bindings: map[reflect.Type]map[string]*binding{}}
}

// bind will map an abstraction to a concrete and set instance if it's a singleton binding.
func (c Container) bind(name string, resolver interface{}, singleton bool) {
	resolverTypeOf := reflect.TypeOf(resolver)
	if resolverTypeOf.Kind() != reflect.Func {
		panic("the resolver must be a function")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for i := 0; i < resolverTypeOf.NumOut(); i++ {
		abstraction := resolverTypeOf.Out(i)
		if _, exists := c.bindings[abstraction][name]; exists {
			panic("concrete already exists  for the abstraction: " + abstraction.String())
		}
		if _, exists := c.bindings[abstraction]; !exists {
			c.bindings[abstraction] = map[string]*binding{}
		}
		c.bindings[abstraction][name] = &binding{
			resolver:  resolver,
			instance:  nil,
			singleton: singleton,
		}
	}
}

// invoke will call the given function and return its returned value.
// It only works for functions that return a single value.
func (c Container) invoke(function interface{}) interface{} {
	return reflect.ValueOf(function).Call(c.arguments("", function))[0].Interface()
}

// arguments will return resolved arguments of the given function.
func (c Container) arguments(name string, function interface{}) []reflect.Value {
	functionTypeOf := reflect.TypeOf(function)
	argumentsCount := functionTypeOf.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := functionTypeOf.In(i)
		arguments[i] = reflect.ValueOf(c.resolve(name, abstraction))
	}

	return arguments
}

func (c Container) resolve(name string, abstraction reflect.Type) interface{} {
	if instance := c.resolveLocally(name, abstraction); instance != nil {
		return instance
	}
	if c.parent != nil {
		return c.parent.resolve(name, abstraction)
	}
	panic("no concrete found for the abstraction: " + abstraction.String())
}

func (c Container) resolveLocally(name string, abstraction reflect.Type) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if binding, ok := c.bindings[abstraction][name]; ok {
		if binding.singleton {
			binding.mu.Lock()
			defer binding.mu.Unlock()

			if binding.instance == nil {
				binding.instance = c.invoke(binding.resolver)
			}
			return binding.instance
		}
		return c.invoke(binding.resolver)
	}
	return nil
}

// Singleton will bind an abstraction to a concrete for further singleton resolves.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have bound already in Container.
func (c Container) Singleton(resolver interface{}) {
	c.SingletonNamed("", resolver)
}

// SingletonNamed will bind a named abstraction to a concrete for further singleton resolves.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have bound already in Container.
func (c Container) SingletonNamed(name string, resolver interface{}) {
	c.bind(name, resolver, true)
}

// Transient will bind an abstraction to a concrete for further transient resolves.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have bound already in Container.
func (c Container) Transient(resolver interface{}) {
	c.TransientNamed("", resolver)
}

// TransientNamed will bind a named abstraction to a concrete for further transient resolves.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have bound already in Container.
func (c Container) TransientNamed(name string, resolver interface{}) {
	c.bind(name, resolver, false)
}

// Reset will reset the container and remove all the bindings.
func (c Container) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k := range c.bindings {
		delete(c.bindings, k)
	}
}

// Make will resolve the dependency and return an appropriate concrete of the given abstraction.
// It can take an abstraction (interface reference) and fill it with the related implementation.
// It also can takes a function (receiver) with one or more arguments of the abstractions (interfaces) that need to be
// resolved, Container will invoke the receiver function and pass the related implementations.
func (c Container) Make(receiver interface{}) {
	c.MakeNamed("", receiver)
}

// MakeNamed will resolve the named dependency and return an appropriate concrete of the given abstraction.
// It can take an abstraction (interface reference) and fill it with the related implementation.
// It also can takes a function (receiver) with one or more arguments of the abstractions (interfaces) that need to be
// resolved, Container will invoke the receiver function and pass the related implementations.
func (c Container) MakeNamed(name string, receiver interface{}) {
	receiverTypeOf := reflect.TypeOf(receiver)
	if receiverTypeOf == nil {
		panic("cannot detect type of the receiver, make sure your are passing reference of the object")
	}

	if receiverTypeOf.Kind() == reflect.Ptr {
		abstraction := receiverTypeOf.Elem()

		instance := c.resolve(name, abstraction)
		reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(instance))
		return
	}

	if receiverTypeOf.Kind() == reflect.Func {
		arguments := c.arguments(name, receiver)
		reflect.ValueOf(receiver).Call(arguments)
		return
	}

	panic("the receiver must be either a reference or a callback")
}

// ForEachNamed iterates over all named concretes
func (c Container) ForEachNamed(function interface{}) {
	functionTypeOf := reflect.TypeOf(function)
	if functionTypeOf.Kind() != reflect.Func {
		panic("argument must be a function")
	}
	if functionTypeOf.NumIn() != 1 {
		panic("function have to accept exactly one argument")
	}
	if functionTypeOf.NumOut() != 0 {
		panic("function must not return anything")
	}
	abstraction := functionTypeOf.In(0)

	c.mu.RLock()
	defer c.mu.RUnlock()
	for name := range c.bindings[abstraction] {
		if name == "" {
			continue
		}
		arguments := c.arguments(name, function)
		reflect.ValueOf(function).Call(arguments)
	}
}

// SubContainer creates sub container
// Bindings are resolved in sub container first and if it's not possible request is redirected to the parent one
func (c Container) SubContainer() Container {
	sub := NewContainer()
	sub.parent = &c
	return sub
}
