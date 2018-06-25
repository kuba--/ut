package ut

import (
	"fmt"
	"testing"
)

func TestTokenize(t *testing.T) {
	x, y := "p(f(X),h(Y,f(a)),Y)", "p(X,h(Z,Y),f(a))"
	n, ix, iy := 15, 0, 8
	tokens := Tokenize(x, y)

	if len(tokens) != n {
		t.Fatalf("len(tokens) Got: %d, Expected: %d", len(tokens), n)
	}

	if x != tokens[ix].Term {
		t.Fatalf("tokens[%d] Got: %s, Expected: %s", ix, tokens[ix].Term, x)
	}

	if y != tokens[iy].Term {
		t.Fatalf("tokens[%d] Got: %s, Expected: %s", iy, tokens[iy].Term, y)
	}
}

func TestUT(t *testing.T) {
	x, y := "p(Z,h(Z,W),f(W))", "p(f(X),h(Y,f(a)),Y)"
	n, ix, iy := 12, 11, 6

	ut := New(Tokenize(x, y))
	if len(ut.Entries) != n {
		for i, e := range ut.Entries {
			t.Logf("[%d]: %s\n", i, e.Term)
		}
		t.Fatalf("len(ut.Entries) Got: %d, Expected: %d", len(ut.Entries), n)
	}

	if ut.Lookup[x] != ix {
		t.Fatalf("ut.Lookup[%s] Got: %v, Expected: %d", x, ut.Lookup[x], ix)
	}

	if ut.Lookup[y] != iy {
		t.Fatalf("ut.Lookup[%s] Got: %v, Expected: %d", y, ut.Lookup[y], iy)
	}
}

func TestUnify(t *testing.T) {
	x, y := "p(Z,h(Z,W),f(W))", "p(f(X),h(Y,f(a)),Y)"
	ut := New(Tokenize(x, y))
	ix, iy := ut.Lookup[x], ut.Lookup[y]
	if !ut.Unify(ix, iy) {
		t.Fatalf("ut.Unify(%d, %d) failed", ix, iy)
	}

	for v1, v2 := range ut.Bindings {
		v2 = ut.dereference(v2)
		t.Logf("%s = %s\n", ut.Entries[v1].Term, ut.TermString(v2))
	}

	mguW := ut.MGU("W")
	if mguW != "f(a)" {
		t.Fatalf("Got W => %s Expected: f(a)", mguW)
	}

	mguX := ut.MGU("X")
	if mguX != "f(a)" {
		t.Fatalf("Got X => %s Expected: f(a)", mguX)
	}

	mguY := ut.MGU("Y")
	if mguY != "f(f(a))" {
		t.Fatalf("Got Y => %s Expected: f(f(a))", mguY)
	}

	mguZ := ut.MGU("Z")
	if mguZ != "f(f(a))" {
		t.Fatalf("Got Z => %s Expected: f(f(a))", mguZ)
	}
}

func ExampleMGU() {
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
}
