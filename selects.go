package treepath

import "github.com/mmbros/demo/queue"

// ----------------------------------------------------------------------------

// selectSelf selects the current element into the candidate list.
type selectSelf struct{}

func (s *selectSelf) apply(e Element, p *pather) {
	p.candidates = append(p.candidates, e)
}

// ----------------------------------------------------------------------------

// selectParent selects the element's parent into the candidate list.
type selectParent struct{}

func (s *selectParent) apply(e Element, p *pather) {
	if parent := e.Parent(); parent != nil {
		p.candidates = append(p.candidates, parent)
	}
}

// ----------------------------------------------------------------------------

// selectChildren selects the element's child elements into the
// candidate list.
type selectChildren struct{}

func (s *selectChildren) apply(e Element, p *pather) {
	for _, child := range e.Children() {
		p.candidates = append(p.candidates, child)
	}
}

// ----------------------------------------------------------------------------

// selectDescendants selects all descendant child elements
// of the element into the candidate list.
type selectDescendants struct{}

func (s *selectDescendants) apply(e Element, p *pather) {
	q := queue.NewFifo(0)

	for q.Push(e); q.Len() > 0; {
		e := q.Pop().(Element)
		p.candidates = append(p.candidates, e)
		for _, c := range e.Children() {
			q.Push(c)
		}
	}
}

// ----------------------------------------------------------------------------

// selectChildrenByTag selects into the candidate list all child
// elements of the element having the specified tag.
type selectChildrenByTag struct {
	tag string
}

func newSelectChildrenByTag(tag string) *selectChildrenByTag {
	return &selectChildrenByTag{tag}
}

func (s *selectChildrenByTag) apply(e Element, p *pather) {
	for _, c := range e.Children() {
		if c.MatchTag(s.tag) {
			p.candidates = append(p.candidates, c)
		}
	}
}
