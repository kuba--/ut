# ut
Package ut implements "Yet Another Efficient Unification Algorithm" by Alin Suciu (https://arxiv.org/abs/cs/0603080v1).

The unification algorithm is at the core of the logic programming paradigm, the first unification algorithm being developed by Robinson. More efficient algorithms were developed later by Martelli and, Montanari.
Here yet another efficient unification algorithm centered on a specific data structure, called the Unification Table.

```Go
	x, y := "p(Z,h(Z,W),f(W))", "p(f(X),h(Y,f(a)),Y)"

	// New Unification Table
	ut := New(Tokenize(x, y))
	ix, iy := ut.Lookup[x], ut.Lookup[y]

	if ut.Unify(ix, iy) {
		fmt.Println("W = " + ut.MGU("W"))
		fmt.Println("X = " + ut.MGU("X"))
		fmt.Println("Y = " + ut.MGU("Y"))
		fmt.Println("Z = " + ut.MGU("Z"))
	}

	// Output:
	// W = f(a)
	// X = f(a)
	// Y = f(f(a))
	// Z = f(f(a))
```
