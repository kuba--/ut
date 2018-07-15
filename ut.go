package ut

import (
	"fmt"
	"log"
	"strings"
)

var (
	// VAR stands for variable type
	VAR = func(t rune) bool {
		return t == Variable
	}

	// STR stands for constants or composite terms
	STR = func(t rune) bool {
		return t == Atom ||
			t == Float ||
			t == Int ||
			t == String ||
			t == Void ||
			t == Functor
	}
)

type (
	// Entry stands for Unification Table Entry
	Entry struct {
		Term       string
		Functor    string
		Components []int
		Type       rune
	}

	// UT stands for Unification Table
	UT struct {
		// Lookup table (term) -> (index)
		Lookup   map[string]int
		Entries  []*Entry
		Bindings map[int]int
	}
)

// Arity is the arity of the term; for variables and constants, it is 0.
func (e *Entry) Arity() int {
	if e.Type != Functor {
		return 0
	}

	return len(e.Components)
}

func (e *Entry) String() string {
	return fmt.Sprintf("{Type: %s, Functor: %s, Components: %v}", TokenString(e.Type), e.Functor, e.Components)
}

// Unify returns a unification maps with VAR bindings.
// Also see ut.MGU for particular terms.
func Unify(x, y string) map[string]string {
	tokens := Tokenize(x, y)
	ut := New(tokens)
	ix, iy := ut.Lookup[x], ut.Lookup[y]
	if !ut.Unify(ix, iy) {
		return nil
	}

	mgu := make(map[string]string)
	for i, j := range ut.Bindings {
		j = ut.Dereference(j)
		mgu[ut.Entries[i].Term] = ut.TermString(j)
	}

	return mgu
}

// New creates a new Unification Table.
func New(tokens []*Token) (ut *UT) {
	ut = &UT{Lookup: make(map[string]int), Bindings: make(map[int]int)}

	for i, n := 0, len(tokens)-1; n >= 0; n-- {
		t := tokens[n]
		if _, exists := ut.Lookup[t.Term]; exists {
			continue
		}

		e := &Entry{
			Term:    t.Term,
			Functor: t.Functor,
			Type:    t.Type,
		}

		for _, c := range t.Components {
			idx, exists := ut.Lookup[c]
			if !exists {
				log.Fatalf("Component: %s not found", c)
			}
			e.Components = append(e.Components, idx)
		}

		ut.Lookup[t.Term] = i
		ut.Entries = append(ut.Entries, e)
		i++
	}

	return ut
}

// Unify tries to calculate MGU (Most General Unifier)
func (ut *UT) Unify(ix, iy int) bool {
	var (
		// stacks
		sx, sy = []int{ix}, []int{iy}
		empty  = func(s []int) bool {
			return len(s) == 0
		}
		pop = func(s []int) (int, []int) {
			n := len(s)
			if n > 0 {
				i := s[n-1]
				s = s[:n-1]

				return i, s
			}

			return 0, nil
		}

		push = func(s []int, i ...int) []int {
			return append(s, i...)
		}
	)

	for !empty(sx) && !empty(sy) {
		// pop entries from stacks
		ix, sx = pop(sx)
		iy, sy = pop(sy)
		ex, ey := ut.Entries[ix], ut.Entries[iy]

		switch true {
		// case 1: ex is bound to a term and ey is bound to a term}
		case STR(ex.Type) && STR(ey.Type):
			ax, ay := ex.Arity(), ey.Arity()
			if ex.Functor != ey.Functor || ax != ay {
				return false
			}

			if ax > 0 {
				sx = push(sx, ex.Components...)
				sy = push(sy, ey.Components...)
			}

		// case 2: ex is bound to a term and ey is bound to a variable
		case STR(ex.Type) && VAR(ey.Type):
			if idx, _, b := ut.bindSTR(ix, iy); !b {
				sx = push(sx, ix)
				sy = push(sy, idx)
			}

		// case 3: ex is bound to a variable and ey is bound to a term
		case VAR(ex.Type) && STR(ey.Type):
			if idx, _, b := ut.bindSTR(iy, ix); !b {
				sx = push(sx, idx)
				sy = push(sy, iy)
			}

		// case 4: ex is bound to a variable and ey is bound to a variable
		case VAR(ex.Type) && VAR(ey.Type):
			if idx1, idx2, b := ut.bindVAR(ix, iy); !b {
				sx = push(sx, idx1)
				sy = push(sy, idx2)
			}
		}
	}

	return true
}

// MGU returns The Most General Unifier as a string for a given term.
// It dereferences bindings and term componenets.
func (ut *UT) MGU(term string) string {
	i, ok := ut.Lookup[term]
	if !ok {
		return ""
	}

	i = ut.Dereference(ut.Bindings[i])
	return ut.TermString(i)
}

// TermString constructs a new term string by dereferencing all components.
func (ut *UT) TermString(idx int) string {
	n := len(ut.Entries)
	if idx < n {
		e := ut.Entries[idx]
		if e.Type != Functor {
			return e.Functor
		}

		components := []string{}
		for _, c := range e.Components {
			i := ut.Dereference(c)
			if i != idx && ut.Entries[i].Type == Functor {
				components = append(components, ut.TermString(i))
			} else {
				components = append(components, ut.Entries[i].Functor)
			}
		}

		return ut.Entries[idx].Functor + "(" + strings.Join(components, ",") + ")"
	}
	return ""
}

// Dereference follows bindings and returns index for dereferenced variable.
func (ut *UT) Dereference(idx int) int {
	i, ok := idx, true
	for ok {
		i, ok = ut.Bindings[i]
		if ok {
			idx = i
		}
	}
	return idx
}

// bindSTR tries to bind a VAR(varIdx) to STR(strIdx).
// If VAR(varIdx) is already bound then dereference it and try again
// or returns indexes to push them on stacks
func (ut *UT) bindSTR(strIdx, varIdx int) (int, int, bool) {
	idx, ok := ut.Bindings[varIdx]
	if !ok {
		// var is a free variable
		// bind var to str
		ut.Bindings[varIdx] = strIdx
		return strIdx, varIdx, true
	}

	// var is already bound - dereference
	idx = ut.Dereference(idx)

	// var is bound to a STR
	e := ut.Entries[idx]
	if STR(e.Type) {
		return idx, varIdx, false
	}

	// free variable
	ut.Bindings[varIdx] = idx
	return idx, varIdx, true
}

// bindVAR tries to bind two VARs.
// If both are already bound the function returns indexes to push them on stacks
// and false as boolean information that binding failed.
func (ut *UT) bindVAR(varIdx1, varIdx2 int) (int, int, bool) {
	i1, ok1 := ut.Bindings[varIdx1]
	i2, ok2 := ut.Bindings[varIdx2]

	// var1 is free and var2 is free
	if !ok1 && !ok2 {
		ut.Bindings[varIdx1] = varIdx2
		return varIdx1, varIdx2, true
	}

	// var1 is free and var2 is bound
	if !ok1 && ok2 {
		ut.Bindings[varIdx1] = varIdx2
		return varIdx1, varIdx2, true
	}

	// var1 is bound and var2 us free
	if ok1 && !ok2 {
		ut.Bindings[varIdx2] = varIdx1
		return varIdx1, varIdx2, true
	}

	// var1 is bound and var2 is bound
	return i1, i2, false
}
