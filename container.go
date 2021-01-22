// Package container provides an IoC container for Go projects.
// It provides simple, fluent and easy-to-use interface to make dependency injection in GoLang easier.
package container

import (
	internal "github.com/golobby/container/pkg/container"
)

func NewContainer() internal.Container {
	return internal.NewContainer()
}

// A default instance for container
var container = NewContainer()

// Singleton creates a singleton for the default instance.
func Singleton(resolver interface{}) {
	container.Singleton(resolver)
}

// SingletonNamed creates a named singleton for the default instance.
func SingletonNamed(name string, resolver interface{}) {
	container.SingletonNamed(name, resolver)
}

// Transient creates a transient binding for the default instance.
func Transient(resolver interface{}) {
	container.Transient(resolver)
}

// TransientNamed creates a named transient binding for the default instance.
func TransientNamed(name string, resolver interface{}) {
	container.TransientNamed(name, resolver)
}

// ForEachNamed iterates over all named concretes
func ForEachNamed(function interface{}) {
	container.ForEachNamed(function)
}

// Reset removes all bindings in the default instance.
func Reset() {
	container.Reset()
}

// Make binds receiver to the default instance.
func Make(receiver interface{}) {
	container.Make(receiver)
}
