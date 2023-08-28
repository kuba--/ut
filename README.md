[![GoDoc](https://godoc.org/github.com/kuba--/ut?status.svg)](http://godoc.org/github.com/kuba--/ut)
[![Go Report Card](https://goreportcard.com/badge/github.com/kuba--/ut)](https://goreportcard.com/report/github.com/kuba--/ut)
[![Build Status](https://github.com/kuba--/ut/workflows/build/badge.svg)](https://github.com/kuba--/ut/actions?query=workflow%3Abuild)
[![Version](https://badge.fury.io/gh/kuba--%2Fut.svg)](https://github.com/kuba--/ut/releases)

# ut
Package ut implements "Yet Another Efficient Unification Algorithm" by Alin Suciu (https://arxiv.org/abs/cs/0603080v1).

The unification algorithm is at the core of the logic programming paradigm, the first unification algorithm being developed by Robinson. More efficient algorithms were developed later by Martelli and, Montanari.

Unification unifies terms. A term is one of variable, constant, functor, or literal value. A variable starts with an upper case letter and their names are the keys in the map returned by Unify. Constants start with a lowercase letter. Values are integers, floating point values, or strings. As a simple example:
```go
x,y := `f(X,Y)`, `f(1,"hello")`
```
This would unify the functor f such that the variable `X` is `1` and the variable `Y` is `"hello"` in the returned map.


Here yet another efficient unification algorithm centered on a specific data structure, called the Unification Table.

```Go
x, y := "p(Z,h(Z,W),f(W))", "p(f(X),h(Y,f(a)),Y)"
mgu := ut.Unify(x, y)
fmt.Println("W = " + mgu["W"])
fmt.Println("X = " + mgu["X"])
fmt.Println("Y = " + mgu["Y"])
fmt.Println("Z = " + mgu["Z"])

// Output:
// W = f(a)
// X = f(a)
// Y = f(f(a))
// Z = f(f(a))

x, y = "f(X1,g(X2,X3),X2,b)", "f(g(h(a,X5),X2),X1,h(a,X4),X4)"
mgu = ut.Unify(x, y)
fmt.Println("X1 = " + mgu["X1"])
fmt.Println("X2 = " + mgu["X2"])
fmt.Println("X3 = " + mgu["X3"])
fmt.Println("X4 = " + mgu["X4"])
fmt.Println("X5 = " + mgu["X5"])

// Output:
// X1 = g(h(a,b),h(a,b))
// X2 = h(a,b)
// X3 = h(a,b)
// X4 = b
// X5 = b
```
