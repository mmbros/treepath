// Credits: https://github.com/beevik/etree

// Package treepath implements the selection of nodes in an arbitrary tree
// of objects with XPath-like espressions.
// See path_test.go for usage example.
package treepath

import (
	"strconv"
	"strings"
)

// Element is the interface that must be satifsfied by a tree node in order to
// enable the treepath FindElements search.
type Element interface {
	// Parent returns the parent of the current element.
	Parent() Element

	// Children returns the children of the current element.
	Children() []Element

	// MatchTag returns true if the current element matches the Tag.
	MatchTag(string) bool

	// MatchTagText returns true if the current element matches both Tag and
	// Text value.
	MatchTagText(string, string) bool

	// MatchAttr returns true if the current element has an attribute that
	// matched the given Name.
	MatchAttr(string) bool

	// MatchAttrText returns true if the current element has an attribute that
	// matched the given Name and Text value.
	MatchAttrText(string, string) bool
}

// Path represents the compiled version of an XPath-like espression.
type Path struct {
	segments []segment
}

// CompilePath creates an optimized version of an XPath-like string that
// can be used to query elements in an element tree.
func CompilePath(path string) (Path, error) {
	var comp compiler
	segments := comp.parsePath(path)
	if comp.err != ErrPath("") {
		return Path{nil}, comp.err
	}
	return Path{segments}, nil
}

// FindElements returns the descendant of root Element that matched the path.
// NOTE: The root element is never matched.
func (path Path) FindElements(root Element) []Element {
	p := newPather()
	return p.traverse(root, path)
}

// A segment is a portion of a path between "/" characters.
// It contains one selector and zero or more [filters].
type segment struct {
	sel     selector
	filters []filter
}

// A selector selects XML elements for consideration by the
// path traversal.
type selector interface {
	apply(e Element, p *pather)
}

// A filter pares down a list of candidate XML elements based
// on a path filter in [brackets].
type filter interface {
	apply(p *pather)
}

// A pather is helper object that traverses an element tree using
// a Path object.  It collects and deduplicates all elements matching
// the path query.
type pather struct {
	queue      *Fifo
	results    []Element
	inResults  map[Element]bool
	candidates []Element
	scratch    []Element // used by filters
}

// A node represents an element and the remaining path segments that
// should be applied against it by the pather.
type node struct {
	e        Element
	segments []segment
}

// A compiler generates a compiled path from a path string.
type compiler struct {
	err ErrPath
}

// ErrPath is returned by path functions when an invalid path is provided.
type ErrPath string

// Error returns the string describing a path error.
func (err ErrPath) Error() string {
	return "treepath: " + string(err)
}

// parsePath parses an XPath-like string describing a path
// through an element tree and returns a slice of segment
// descriptors.
func (c *compiler) parsePath(path string) []segment {
	// If path starts or ends with //, fix it
	if strings.HasPrefix(path, "//") {
		path = "." + path
	}
	if strings.HasSuffix(path, "//") {
		path = path + "*"
	}

	// Paths cannot be absolute
	if strings.HasPrefix(path, "/") {
		c.err = ErrPath("paths cannot be absolute.")
		return nil
	}

	// Split path into segment objects
	var segments []segment
	for _, s := range splitPath(path) {
		segments = append(segments, c.parseSegment(s))
		if c.err != ErrPath("") {
			break
		}
	}
	return segments
}

// splitPath splits a path in the segments between / characters.
// It handles the / characters eventually contained in the text values of path.
func splitPath(path string) []string {
	pieces := make([]string, 0)
	start := 0
	inquote := false
	for i := 0; i+1 <= len(path); i++ {
		if path[i] == '\'' {
			inquote = !inquote
		} else if path[i] == '/' && !inquote {
			pieces = append(pieces, path[start:i])
			start = i + 1
		}
	}
	return append(pieces, path[start:])
}

// parseSegment parses a path segment between / characters.
func (c *compiler) parseSegment(path string) segment {
	pieces := strings.Split(path, "[")
	seg := segment{
		sel:     c.parseSelector(pieces[0]),
		filters: make([]filter, 0),
	}
	for i := 1; i < len(pieces); i++ {
		fpath := pieces[i]
		last := len(fpath) - 1
		if last < 0 || fpath[last] != ']' {
			c.err = ErrPath("path has invalid filter [brackets].")
			break
		}
		seg.filters = append(seg.filters, c.parseFilter(fpath[:last]))
	}
	return seg
}

// parseSelector parses a selector at the start of a path segment.
func (c *compiler) parseSelector(path string) selector {
	switch path {
	case ".":
		return new(selectSelf)
	case "..":
		return new(selectParent)
	case "*":
		return new(selectChildren)
	case "":
		return new(selectDescendants)
	default:
		return newSelectChildrenByTag(path)
	}
}

// parseFilter parses a path filter contained within [brackets].
func (c *compiler) parseFilter(path string) filter {

	if len(path) == 0 {
		c.err = ErrPath("path contains an empty filter expression.")
		return nil
	}

	// Filter contains [@attr='text'] or [tag='text']?
	eqindex := strings.Index(path, "='")
	if eqindex >= 0 {
		rindex := nextIndex(path, "'", eqindex+2)
		if rindex != len(path)-1 {
			c.err = ErrPath("path has mismatched filter quotes.")
			return nil
		}
		switch {
		case path[0] == '@':
			// [@attr='text']
			return newFilterAttrText(path[1:eqindex], path[eqindex+2:rindex])
		default:
			// [tag='text']
			return newFilterChildText(path[:eqindex], path[eqindex+2:rindex])
		}
	}
	// Filter contains [@attr], [N] or [tag]
	switch {
	case path[0] == '@':
		// [@attr]
		return newFilterAttr(path[1:])
	case isInteger(path):
		// [N]
		pos, _ := strconv.Atoi(path)
		switch {
		case pos > 0:
			return newFilterPos(pos - 1)
		default:
			return newFilterPos(pos)
		}
	default:
		// [tag]
		return newFilterChild(path)
	}

}

func (seg *segment) apply(e Element, p *pather) {
	seg.sel.apply(e, p)
	for _, f := range seg.filters {
		f.apply(p)
	}
}

func newPather() *pather {
	return &pather{
		queue:      NewFifo(0),
		results:    make([]Element, 0),
		inResults:  make(map[Element]bool),
		candidates: make([]Element, 0),
		scratch:    make([]Element, 0),
	}
}

// traverse follows the path from the element e, collecting
// and then returning all elements that match the path's selectors
// and filters.
func (p *pather) traverse(e Element, path Path) []Element {
	for p.queue.Push(node{e, path.segments}); p.queue.Len() > 0; {
		p.eval(p.queue.Pop().(node))
	}
	return p.results
}

// eval evalutes the current path node by applying the remaining
// path's selector rules against the node's element.
func (p *pather) eval(n node) {
	p.candidates = p.candidates[0:0]
	seg, remain := n.segments[0], n.segments[1:]
	seg.apply(n.e, p)

	if len(remain) == 0 {
		for _, c := range p.candidates {
			if in := p.inResults[c]; !in {
				p.inResults[c] = true
				p.results = append(p.results, c)
			}
		}
	} else {
		for _, c := range p.candidates {
			p.queue.Push(node{c, remain})
		}
	}
}
