package treepath

import "testing"

type myNode int

func TestFifo(t *testing.T) {
	const L = 20

	// init fifo queue with 0 size
	q := NewFifo(0)
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
		q.Push(myNode(j))
		if q.Len() != j {
			t.Errorf("Len: expected %d, found %d", j, q.Len())
		}
	}

	// pop L elements and test order
	j := 1
	for q.Len() > 0 {
		mynode := q.Pop().(myNode)
		n := int(mynode)
		if j != n {
			t.Errorf("Pop: expected %d, found %d", j, n)

		}
		j++
		if j > 1000 {
			t.Errorf("Pop: too many iterations")
		}
	}

}
