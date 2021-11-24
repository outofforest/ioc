package ioc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Shape interface {
	SetArea(int)
	GetArea() int
}

type Circle struct {
	a int
}

func (c *Circle) SetArea(a int) {
	c.a = a
}

func (c Circle) GetArea() int {
	return c.a
}

type Database interface {
	Connect() bool
}

type MySQL struct{}

func (m MySQL) Connect() bool {
	return true
}

type Concrete interface {
	String() string
}

type ConcreteA struct {
}

func (c ConcreteA) String() string {
	return "A"
}

type ConcreteB struct {
}

func (c ConcreteB) String() string {
	return "B"
}

type ConcreteC struct {
}

func (c ConcreteC) String() string {
	return "C"
}

func TestSingletonItShouldMakeAnInstanceOfTheAbstraction(t *testing.T) {
	area := 5

	instance := New()
	instance.Singleton(func() Shape {
		return &Circle{a: area}
	})

	instance.Resolve(func(s Shape) {
		a := s.GetArea()
		assert.Equal(t, area, a)
	})
}

func TestSingletonItShouldMakeSameObjectEachMake(t *testing.T) {
	instance := New()
	instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})

	area := 6

	instance.Resolve(func(s1 Shape) {
		s1.SetArea(area)
	})

	instance.Resolve(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, area)
	})
}

func TestSingletonNamedInstanceIsDifferentThanDefaultOne(t *testing.T) {
	instance := New()
	instance.Singleton(func() Concrete {
		return &ConcreteA{}
	})
	instance.SingletonNamed("named", func() Concrete {
		return &ConcreteB{}
	})

	var concreteA Concrete
	var concreteB Concrete

	instance.Resolve(&concreteA)
	instance.ResolveNamed("named", &concreteB)

	assert.Equal(t, "A", concreteA.String())
	assert.Equal(t, "B", concreteB.String())
}

func TestSingletonTwoNamedInstancesAreDifferent(t *testing.T) {
	instance := New()
	instance.SingletonNamed("A", func() Concrete {
		return &ConcreteA{}
	})
	instance.SingletonNamed("B", func() Concrete {
		return &ConcreteB{}
	})

	var concreteA Concrete
	var concreteB Concrete

	instance.ResolveNamed("A", &concreteA)
	instance.ResolveNamed("B", &concreteB)

	assert.Equal(t, "A", concreteA.String())
	assert.Equal(t, "B", concreteB.String())
}

func TestSingletonWithNonFunctionResolverItShouldPanic(t *testing.T) {
	instance := New()
	value := "the resolver must be a function"
	assert.PanicsWithValue(t, value, func() {
		instance.Singleton("STRING!")
	}, "Expected panic")
}

func TestSingletonItShouldResolveResolverArguments(t *testing.T) {
	area := 5
	instance := New()
	instance.Singleton(func() Shape {
		return &Circle{a: area}
	})

	instance.Singleton(func(s Shape) Database {
		assert.Equal(t, s.GetArea(), area)
		return &MySQL{}
	})
}

func TestTransientItShouldMakeDifferentObjectsOnMake(t *testing.T) {
	area := 5
	instance := New()
	instance.Transient(func() Shape {
		return &Circle{a: area}
	})

	instance.Resolve(func(s1 Shape) {
		s1.SetArea(6)
	})

	instance.Resolve(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, area)
	})
}

func TestTransientNamedInstanceIsDifferentThanDefaultOne(t *testing.T) {
	instance := New()
	instance.Transient(func() Concrete {
		return &ConcreteA{}
	})
	instance.TransientNamed("named", func() Concrete {
		return &ConcreteB{}
	})

	var concreteA Concrete
	var concreteB Concrete

	instance.Resolve(&concreteA)
	instance.ResolveNamed("named", &concreteB)

	assert.Equal(t, "A", concreteA.String())
	assert.Equal(t, "B", concreteB.String())
}

func TestTransientTwoNamedInstancesAreDifferent(t *testing.T) {
	instance := New()
	instance.TransientNamed("A", func() Concrete {
		return &ConcreteA{}
	})
	instance.TransientNamed("B", func() Concrete {
		return &ConcreteB{}
	})

	var concreteA Concrete
	var concreteB Concrete

	instance.ResolveNamed("A", &concreteA)
	instance.ResolveNamed("B", &concreteB)

	assert.Equal(t, "A", concreteA.String())
	assert.Equal(t, "B", concreteB.String())
}

func TestTransientItShouldMakeAnInstanceOfTheAbstraction(t *testing.T) {
	area := 5
	instance := New()
	instance.Transient(func() Shape {
		return &Circle{a: area}
	})

	instance.Resolve(func(s Shape) {
		a := s.GetArea()
		assert.Equal(t, a, area)
	})
}

func TestResolveWithSingleInputAndCallback(t *testing.T) {
	instance := New()
	instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})

	instance.Resolve(func(s Shape) {
		if _, ok := s.(*Circle); !ok {
			t.Error("Expected Circle")
		}
	})
}

func TestResolveWithMultipleInputsAndCallback(t *testing.T) {
	instance := New()
	instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})

	instance.Singleton(func() Database {
		return &MySQL{}
	})

	instance.Resolve(func(s Shape, m Database) {
		if _, ok := s.(*Circle); !ok {
			t.Error("Expected Circle")
		}

		if _, ok := m.(*MySQL); !ok {
			t.Error("Expected MySQL")
		}
	})
}

func TestResolveWithSingleInputAndReference(t *testing.T) {
	instance := New()
	instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})

	var s Shape

	instance.Resolve(&s)

	if _, ok := s.(*Circle); !ok {
		t.Error("Expected Circle")
	}
}

func TestResolveWithMultipleInputsAndReference(t *testing.T) {
	instance := New()
	instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})

	instance.Singleton(func() Database {
		return &MySQL{}
	})

	var (
		s Shape
		d Database
	)

	instance.Resolve(&s)
	instance.Resolve(&d)

	if _, ok := s.(*Circle); !ok {
		t.Error("Expected Circle")
	}

	if _, ok := d.(*MySQL); !ok {
		t.Error("Expected MySQL")
	}
}

func TestResolveWithUnsupportedReceiver(t *testing.T) {
	instance := New()
	value := "the receiver must be either a reference or a callback"
	assert.PanicsWithValue(t, value, func() {
		instance.Resolve("STRING!")
	}, "Expected panic")
}

func TestResolveWithNonReference(t *testing.T) {
	instance := New()
	value := "cannot detect type of the receiver, make sure your are passing reference of the object"
	assert.PanicsWithValue(t, value, func() {
		var s Shape
		instance.Resolve(s)
	}, "Expected panic")
}

func TestResolveWithUnboundedAbstraction(t *testing.T) {
	instance := New()
	value := "no concrete found for the abstraction: ioc.Shape"
	assert.PanicsWithValue(t, value, func() {
		var s Shape
		instance.Reset()
		instance.Resolve(&s)
	}, "Expected panic")
}

func TestResolveWithCallbackThatHasAUnboundedAbstraction(t *testing.T) {
	instance := New()
	value := "no concrete found for the abstraction: ioc.Database"
	assert.PanicsWithValue(t, value, func() {
		instance.Reset()
		instance.Singleton(func() Shape {
			return &Circle{}
		})
		instance.Resolve(func(s Shape, d Database) {})
	}, "Expected panic")
}

func TestForEachNamedReturnsNamedConcretesOnly(t *testing.T) {
	instance := New()
	instance.Singleton(func() Concrete {
		return &ConcreteA{}
	})
	instance.SingletonNamed("B", func() Concrete {
		return &ConcreteB{}
	})
	instance.TransientNamed("C", func() Concrete {
		return &ConcreteC{}
	})
	instance.SingletonNamed("B", func() Shape {
		return &Circle{}
	})

	vals := map[string]bool{}
	instance.ForEachNamed(func(c Concrete) {
		assert.False(t, vals[c.String()])
		vals[c.String()] = true
	})
	assert.Len(t, vals, 2)
	assert.True(t, vals["B"])
	assert.True(t, vals["C"])
}

func TestForEachNamedResolvesFromParentInstance(t *testing.T) {
	instance := New()
	instance.Singleton(func() Concrete {
		return &ConcreteA{}
	})
	instance.SingletonNamed("B", func() Concrete {
		return &ConcreteB{}
	})
	subInstance := instance.SubContainer()
	subInstance.TransientNamed("C", func() Concrete {
		return &ConcreteC{}
	})
	subInstance.SingletonNamed("B", func() Shape {
		return &Circle{}
	})

	vals := map[string]bool{}
	subInstance.ForEachNamed(func(c Concrete) {
		assert.False(t, vals[c.String()])
		vals[c.String()] = true
	})
	assert.Len(t, vals, 2)
	assert.True(t, vals["B"])
	assert.True(t, vals["C"])
}

func TestNamesReturnsCorrectNames(t *testing.T) {
	instance := New()
	instance.Singleton(func() Concrete {
		return &ConcreteA{}
	})
	instance.SingletonNamed("B", func() Concrete {
		return &ConcreteB{}
	})
	instance.TransientNamed("C", func() Concrete {
		return &ConcreteC{}
	})
	instance.SingletonNamed("B", func() Shape {
		return &Circle{}
	})

	vals := map[string]bool{}
	for _, name := range instance.Names((*Concrete)(nil)) {
		assert.False(t, vals[name])
		vals[name] = true
	}
	assert.Len(t, vals, 2)
	assert.True(t, vals["B"])
	assert.True(t, vals["C"])
}

func TestNamesIncludesNamesFromParentInstance(t *testing.T) {
	instance := New()
	instance.Singleton(func() Concrete {
		return &ConcreteA{}
	})
	instance.SingletonNamed("B", func() Concrete {
		return &ConcreteB{}
	})
	subInstance := instance.SubContainer()
	subInstance.TransientNamed("C", func() Concrete {
		return &ConcreteC{}
	})
	subInstance.SingletonNamed("B", func() Shape {
		return &Circle{}
	})

	vals := map[string]bool{}
	for _, name := range subInstance.Names((*Concrete)(nil)) {
		assert.False(t, vals[name])
		vals[name] = true
	}
	assert.Len(t, vals, 2)
	assert.True(t, vals["B"])
	assert.True(t, vals["C"])
}

func TestNameExists(t *testing.T) {
	instance := New()
	instance.Singleton(func() Concrete {
		return &ConcreteA{}
	})
	instance.SingletonNamed("B", func() Concrete {
		return &ConcreteB{}
	})
	subInstance := instance.SubContainer()
	subInstance.TransientNamed("C", func() Concrete {
		return &ConcreteC{}
	})
	subInstance.SingletonNamed("B", func() Shape {
		return &Circle{}
	})

	assert.False(t, subInstance.NameExists("A", (*Concrete)(nil)))
	assert.True(t, subInstance.NameExists("B", (*Concrete)(nil)))
	assert.True(t, subInstance.NameExists("C", (*Concrete)(nil)))
}

func TestEachContainerResolvesFromItself(t *testing.T) {
	instance := New()
	instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	subInstance := instance.SubContainer()
	subInstance.Singleton(func() Shape {
		return &Circle{a: 6}
	})

	var shape Shape
	instance.Resolve(&shape)
	assert.Equal(t, 5, shape.GetArea())

	subInstance.Resolve(&shape)
	assert.Equal(t, 6, shape.GetArea())
}

func TestSubContainerResolvesFromParent(t *testing.T) {
	instance := New()
	instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	subInstance := instance.SubContainer()

	var shape Shape
	subInstance.Resolve(&shape)
	assert.Equal(t, 5, shape.GetArea())
}

func TestResolvePanicsIfDependencyResolvesToNil(t *testing.T) {
	instance := New()
	instance.Transient(func() Shape {
		return &Circle{a: 5}
	})
	subInstance := instance.SubContainer()
	subInstance.Transient(func() Shape {
		return nil
	})

	value := "no concrete found for the abstraction: ioc.Shape"
	assert.PanicsWithValue(t, value, func() {
		var shape Shape
		subInstance.Resolve(&shape)
		fmt.Println(shape)
	}, "Expected panic")
}

func TestCallReturningError(t *testing.T) {
	instance := New()
	instance.Transient(func() Shape {
		return &Circle{a: 5}
	})

	var err error
	instance.Call(func(shape Shape) error {
		return errors.New("test error")
	}, &err)
	assert.Equal(t, "test error", err.Error())
}

func TestCallReturningNilError(t *testing.T) {
	instance := New()
	instance.Transient(func() Shape {
		return &Circle{a: 5}
	})

	var err error
	instance.Call(func(shape Shape) error {
		return nil
	}, &err)
	assert.Nil(t, err)
}

func TestCallReturningTwoValues(t *testing.T) {
	instance := New()

	var a int
	var b bool
	instance.Call(func() (int, bool) {
		return 1, true
	}, &a, &b)
	assert.Equal(t, 1, a)
	assert.Equal(t, true, b)
}
