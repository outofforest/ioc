# Container
A lightweight yet powerful IoC container for Go projects. It provides a simple, fluent and easy-to-use interface to make dependency injection in GoLang easier.

## Documentation

### Required Go Versions
It requires Go `v1.15` or newer versions.

### Installation
To install this package, run the following command in the root of your project.

```bash
go get github.com/wojciech-malota-wojcik/ioc
```

### Introduction
GoLobby Container like any other IoC container is used to bind abstractions to their implementations.
Binding is a process of introducing an IoC container that which concrete (implementation) is appropriate for an abstraction. In this process, you also determine how it must be resolved, singleton or transient. 
In singleton binding, the container provides an instance once and returns it for each request. 
In transient binding, the container always returns a brand new instance for each request.
After the binding process, you can ask the IoC container to get the appropriate implementation of the abstraction that your code depends on. In this case, your code depends on abstractions, not implementations.

### Binding

#### Singleton

Singleton binding using Container:

```go
c := ioc.New()
c.Singleton(func() Abstraction {
  return Implementation
})
```

It takes a resolver function which its return type is the abstraction and the function body configures the related concrete (implementation) and returns it.

Example for a singleton binding:

```go
c := ioc.New()
c.Singleton(func() Database {
  return &MySQL{}
})
```

#### Transient

Transient binding is also similar to singleton binding.

Example for a transient binding:

```go
c := ioc.New()
c.Transient(func() Shape {
  return &Rectangle{}
})
```

### Resolving

Container resolves the dependencies with the method `make()`.

#### Using References

One way to get the appropriate implementation you need is to declare an instance of the abstraction type and pass its reference to Container this way:

```go
c := ioc.New()
var a Abstraction
c.Resolve(&a)
// "a" will be implementation of the Abstraction
```

Example:

```go
c := ioc.New()
var m Mailer
c.Resolve(&m)
m.Send("info@miladrahimi.com", "Hello Milad!")
```

#### Using Closures

Another way to resolve the dependencies is by using a function (receiver) that its arguments are the abstractions you 
need. Container will invoke the function and pass the related implementations for each abstraction.

```go
c := ioc.New()
c.Resolve(func(a Abstraction) {
  // "a" will be implementation of the Abstraction
})
```

Example:

```go
c := ioc.New()
c.Resolve(func(db Database) {
  // "db" will be the instance of MySQL
  db.Query("...")
})
```

You can also resolve multiple abstractions this way:

```go
c := ioc.New()
c.Resolve(func(db Database, s Shape) {
  db.Query("...")
  s.Area()
})
```

#### Calling functions and getting results

You may call a function and get returned values:

```go
c := ioc.New()
var result int
var err error
c.Call(func(db Database) (int, error) {
  return db.GetIntValue("...")
}, &result, &int)
```

### Named bindings

You can also use named bindings to create many bindings of the same type:

```go
c := ioc.New()
c.SingletonNamed("concreteFactoryA", func() Factory {
	return &ConcreteFactoryA{}
})
c.SingletonNamed("concreteFactoryB", func() Factory {
    return &ConcreteFactoryB{}
})
```

Then you may easily retrieve concrete factory based on a string saved ex. in DB:

```go
c := ioc.New()
var factory Factory
c.ResolveNamed(factoryName, &factory)
```

You can also easily iterate over all named bindings of particular interface:

```go
c.ForEachNamed(func(factory Factory)) {
    
}
```

It's possible to get list off all the names used for named bindings for a type:

```go
names := c.Names(Factory(nil))
```

### Binding time

You can also resolve a dependency at the binding time in your resolver function like the following example.

```go
c := ioc.New()

// Bind Config to JsonConfig
c.Singleton(func() Config {
    return &JsonConfig{...}
})

// Bind Database to MySQL
c.Singleton(func(c Config) Database {
    // "c" will be the instance of JsonConfig
    return &MySQL{
        Username: c.Get("DB_USERNAME"),
        Password: c.Get("DB_PASSWORD"),
    }
})
```

Notice: You can only resolve the dependencies in a binding resolver function that has already bound.

### Sub containers

You may create sub container:

```go
c := ioc.New()
c.Singleton(func() Binding1 { ... })
subC := c.SubContainer()
subC.Singleton(func() Binding2 { ... })
```

In above case these work:

```go
var binding1 Binding1
var binding2 Binding2
c.Resolve(&binding1)
subC.Resolve(&binding1)
subC.Resolve(&binding2)
```

This doesn't work:

```go
var binding2 Binding2
c.Resolve(&binding2)
```

because Binding2 is registered only in sub container.

### Usage Tips

#### Performance
The package Container inevitably uses reflection in binding and resolving processes. 
If performance is a concern, you should use this package more carefully. 
Try to bind and resolve the dependencies out of the processes that are going to run many times 
(for example, on each request), put it where that run only once when you run your applications 
like main and init functions.

## License

GoLobby Container is released under the [MIT License](http://opensource.org/licenses/mit-license.php).
