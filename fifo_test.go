package treepath

import "testing"

type myElement struct {
	number int
}

func (e *myElement) Parent() Element                      { return nil }
func (e *myElement) Children() []Element                  { return nil }
func (e *myElement) MatchTag(tag string) bool             { return false }
func (e *myElement) MatchTagText(tag, text string) bool   { return false }
func (e *myElement) MatchAttr(attr string) bool           { return false }
func (e *myElement) MatchAttrText(attr, text string) bool { return false }

func TestFifo(t *testing.T) {
	const L = 20

	// init fifo queue with 0 size
	q := newFifo()
	if q.Len() != 0 {
		t.Errorf("Len: expected %d, found %d", 0, q.Len())
	}

	// test pop with len = 0
	p := q.Pop()
	if p != nil {
		t.Errorf("Pop: expected nil, found %v", p)
	}
	// push L elements
	for j := 1; j <= L; j++ {

		q.Push(&node{&myElement{j}, nil})
		if q.Len() != j {
			t.Errorf("Len: expected %d, found %d", j, q.Len())
		}
	}

	// pop L elements and test order
	j := 1
	for q.Len() > 0 {
		mynode := q.Pop()
		n := mynode.e.(*myElement).number
		if j != n {
			t.Errorf("Pop: expected %d, found %d", j, n)

		}
		j++
		if j > 1000 {
			t.Errorf("Pop: too many iterations")
		}
	}
}
