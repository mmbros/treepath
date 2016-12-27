package treepath

// ----------------------------------------------------------------------------

// filterAttrVal filters the candidate list for elements having
// the specified attribute with the specified value.
type filterAttrText struct {
	attr, text string
}

func newFilterAttrText(attr, text string) *filterAttrText {
	return &filterAttrText{attr, text}
}

func (f *filterAttrText) apply(p *pather) {
	for _, c := range p.candidates {
		if c.MatchAttrText(f.attr, f.text) {
			p.scratch = append(p.scratch, c)
		}
	}
	p.candidates, p.scratch = p.scratch, p.candidates[0:0]
}

// ----------------------------------------------------------------------------

// filterPos filters the candidate list, keeping only the
// candidate at the specified index.
type filterPos struct {
	index int
}

func newFilterPos(pos int) *filterPos {
	return &filterPos{pos}
}

func (f *filterPos) apply(p *pather) {
	if f.index >= 0 {
		if f.index < len(p.candidates) {
			p.scratch = append(p.scratch, p.candidates[f.index])
		}
	} else {
		if -f.index <= len(p.candidates) {
			p.scratch = append(p.scratch, p.candidates[len(p.candidates)+f.index])
		}
	}
	p.candidates, p.scratch = p.scratch, p.candidates[0:0]
}

// ----------------------------------------------------------------------------

// filterAttr filters the candidate list for elements having
// the specified attribute.
type filterAttr struct {
	attr string
}

func newFilterAttr(attr string) *filterAttr {
	return &filterAttr{attr}
}

func (f *filterAttr) apply(p *pather) {
	for _, c := range p.candidates {
		if c.MatchAttr(f.attr) {
			p.scratch = append(p.scratch, c)
		}
	}
	p.candidates, p.scratch = p.scratch, p.candidates[0:0]
}

// ----------------------------------------------------------------------------

// filterChild filters the candidate list for elements having
// a child element with the specified tag.
type filterChild struct {
	tag string
}

func newFilterChild(tag string) *filterChild {
	return &filterChild{tag}
}

func (f *filterChild) apply(p *pather) {
	for _, c := range p.candidates {
		for _, cc := range c.Children() {
			if cc.MatchTag(f.tag) {
				p.scratch = append(p.scratch, c)
			}
		}
	}
	p.candidates, p.scratch = p.scratch, p.candidates[0:0]
}

// ----------------------------------------------------------------------------

// filterChildText filters the candidate list for elements having
// a child element with the specified tag and text.
type filterChildText struct {
	tag, text string
}

func newFilterChildText(tag, text string) *filterChildText {
	return &filterChildText{tag, text}
}

func (f *filterChildText) apply(p *pather) {
	for _, c := range p.candidates {
		for _, cc := range c.Children() {
			if cc.MatchTagText(f.tag, f.text) {
				p.scratch = append(p.scratch, c)
			}
		}
	}
	p.candidates, p.scratch = p.scratch, p.candidates[0:0]
}
