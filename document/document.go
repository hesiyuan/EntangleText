package document

import (
	"bytes"
	"math/rand"
)

// Adapted from Ravern Koh's implementation
// Document represents a Logoot Documentument. Actions like Insert and Delete can be performed
// on Document. If at any time an invalid position is given, a panic will occur, so raw
// positions should only be used for debugging purposes.
type Document struct {
	clientID uint8
	pairs    []pair
}

// Pos is an element of a position identifier. A position identifier identifies an
// atom within a Doc. The behaviour of an empty position identifier (length == 0) is
// undefined, so just do not pass in empty position identifiers to any method/function.
type Identifier struct {
	Ident uint16
	Site  uint8
}

// pair is a position identifier and its atom.
type pair struct {
	pos  []Identifier // a position is a list of identifiers
	atom string       // this is actually a char, but stick to string for easy future extension to string-wise

}

// Start and end positions. These will always exist within a Documentument.
var (
	Start = []Identifier{{0, 0}}
	End   = []Identifier{{^uint16(0), 0}}
)

// New creates a new Documentument containing the given content.
func New(content []string, clientID uint8) *Document {
	d := &Document{clientID: clientID}
	d.Insert(Start, "")
	d.Insert(End, "")
	for _, c := range content {
		// End will always exist.
		d.InsertLeft(End, c)
	}
	d.clientID = clientID
	return d
}

/* Basic methods */

// Index of a position in the Document. Secondary value indicates whether the value exists.
// If the value doesn't exist, the index returned is the index that the position would
// have been in, should it have existed.
func (d *Document) Index(p []Identifier) (int, bool) {
	off := 0
	pr := d.pairs
	for {
		if len(pr) == 0 {
			return off, false
		}
		spt := len(pr) / 2 // binary search
		pair := pr[spt]
		if cmp := ComparePos(pair.pos, p); cmp == 0 {
			return spt + off, true
		} else if cmp == -1 {
			off += spt + 1
			pr = pr[spt+1:]
		} else if cmp == 1 {
			pr = pr[0:spt]
		}
	}
}

// ComparePos compares two position identifiers, returning -1 if the left is less than the
// right, 0 if equal, and 1 if greater.
func ComparePos(lp []Identifier, rp []Identifier) int8 {
	for i := 0; i < len(lp); i++ {
		if len(rp) == i {
			return 1
		}
		if lp[i].Ident < rp[i].Ident {
			return -1
		}
		if lp[i].Ident > rp[i].Ident {
			return 1
		}
		if lp[i].Site < rp[i].Site {
			return -1
		}
		if lp[i].Site > rp[i].Site {
			return 1
		}
	}
	if len(rp) > len(lp) {
		return -1
	}
	return 0
}

// Atom at the position. Secondary return value indicates whether the value exists.
func (d *Document) Get(p []Identifier) (string, bool) {
	i, exists := d.Index(p)
	if !exists {
		return "", false
	}
	return d.pairs[i].atom, true
}

// Insert a new pair at the position, returning success or failure (already existing
// position).
func (d *Document) Insert(p []Identifier, atom string) bool {
	i, exists := d.Index(p)
	if exists {
		return false
	}
	d.pairs = append(d.pairs[0:i], append([]pair{{p, atom}}, d.pairs[i:]...)...)
	return true
}

// Delete the pair at the position, returning success or failure (non-existent position).
func (d *Document) Delete(p []Identifier) bool {
	i, exists := d.Index(p)
	if !exists || i == 0 || i == len(d.pairs)-1 {
		return false
	}
	d.pairs = append(d.pairs[0:i], d.pairs[i+1:]...)
	return true
}

// Left returns the position to the left of the given position, and a flag indicating
// whether it exists (when the given position is the start, there is no position to the
// left of it). Will be false if the given position is invalid. The Start pair is not
// considered as an actual pair.
func (d *Document) Left(p []Identifier) ([]Identifier, bool) {
	i, exists := d.Index(p)
	if !exists || i == 0 {
		return nil, false
	}
	return d.pairs[i-1].pos, true
}

// Right returns the position to the right of the given position, and a flag indicating
// whether it exists (when the given position is the end, there is no position to the
// right of it). Will be false if the given position is invalid. The End pair is not
// considered as an actual pair.
func (d *Document) Right(p []Identifier) ([]Identifier, bool) {
	i, exists := d.Index(p)
	if !exists || i >= len(d.pairs)-1 {
		return nil, false
	}
	return d.pairs[i+1].pos, true
}

// random number between x and y, where y is greater than x.
func random(x, y uint16) uint16 {
	return uint16(rand.Intn(int(y-x-1))) + 1 + x
}

// GeneratePos generates a new position identifier between the two positions provided.
// Secondary return value indicates whether it was successful (when the two positions
// are equal, or the left is greater than right, position cannot be generated).
// Later will be optimized if time permits to LSEQ
func GeneratePos(lp, rp []Identifier, site uint8) ([]Identifier, bool) {
	if ComparePos(lp, rp) != -1 {
		return nil, false
	}
	p := []Identifier{}
	for i := 0; i < len(lp); i++ {
		l := lp[i]
		r := rp[i]
		if l.Ident == r.Ident && l.Site == r.Site {
			p = append(p, Identifier{l.Ident, l.Site})
			continue
		}
		if d := r.Ident - l.Ident; d > 1 {
			r := random(l.Ident, r.Ident)
			p = append(p, Identifier{r, site})
		} else if d == 1 {
			if site > l.Site {
				p = append(p, Identifier{l.Ident, site})
			} else if site < r.Site {
				p = append(p, Identifier{r.Ident, site})
			} else {
				min := uint16(0)
				if len(lp) > len(rp) {
					min = lp[len(rp)].Ident
					// Super edge case
					// left  => {3 1} {65534 1}
					// right => {4 1}
					// In this case, 65534 can't be min, because no number is in between
					// it and MAX. So need to extend the positions further.
					if min == ^uint16(0)-1 {
						r := random(0, ^uint16(0))
						p = append(p, Identifier{l.Ident, l.Site})
						p = append(p, lp[len(rp):]...)
						p = append(p, Identifier{r, site})
						return p, true
					}
				}
				r := random(min, ^uint16(0))
				p = append(p, Identifier{l.Ident, l.Site}, Identifier{r, site})
			}
		} else {
			if site > l.Site && site < r.Site {
				p = append(p, Identifier{l.Ident, site})
			} else {
				r := random(0, ^uint16(0))
				p = append(p, Identifier{l.Ident, l.Site}, Identifier{r, site})
			}
		}
		return p, true
	}
	if len(rp) > len(lp) {
		r := random(0, rp[len(lp)].Ident)
		p = append(p, Identifier{r, site})
	}
	return p, true
}

// GeneratePos generates a new position identifier between the two positions provided.
// Secondary return value indicates whether it was successful (when the two positions
// are equal, or the left is greater than right, position cannot be generated).
func (d *Document) GeneratePos(lp []Identifier, rp []Identifier) ([]Identifier, bool) {
	return GeneratePos(lp, rp, d.clientID)
}

/* Convenience methods */

// InsertLeft inserts the atom to the left of the given position, returning the inserted
// position and whether it is successful (when the given position doesn't exist,
// InsertLeft won't do anything and return false).
func (d *Document) InsertLeft(p []Identifier, atom string) ([]Identifier, bool) {
	lp, success := d.Left(p)
	if !success {
		return nil, false
	}
	np, success := d.GeneratePos(lp, p)
	if !success {
		return nil, false
	}
	return np, d.Insert(np, atom)
}

// InsertRight inserts the atom to the right of the given position, returning the inserted
// position whether it is successful (when the given position doesn't exist, InsertRight
// won't do anything and return false).
func (d *Document) InsertRight(p []Identifier, atom string) ([]Identifier, bool) {
	rp, success := d.Right(p)
	if !success {
		return nil, false
	}
	np, success := d.GeneratePos(p, rp)
	if !success {
		return nil, false
	}
	return np, d.Insert(np, atom)
}

// DeleteLeft deletes the atom to the left of the given position, returning whether it
// was successful (when the given position is the start, there is no position to the left
// of it).
func (d *Document) DeleteLeft(p []Identifier) bool {
	lp, success := d.Left(p)
	if !success {
		return false
	}
	return d.Delete(lp)
}

// DeleteRight deletes the atom to the right of the given position, returning whether it
// was successful (when the given position is the end, there is no position to the right
// of it).
func (d *Document) DeleteRight(p []Identifier) bool {
	rp, success := d.Right(p)
	if !success {
		return false
	}
	return d.Delete(rp)
}

// Content of the entire Documentument.
func (d *Document) Content() string {
	var b bytes.Buffer
	for i := 1; i < len(d.pairs)-1; i++ {
		b.WriteString(d.pairs[i].atom)
	}
	return b.String()
}

// Other useful functions for serialization

// PosBytes returns the position as a byte slice.
func PosBytes(p []Identifier) []byte {
	b := []byte{byte(len(p))}
	for _, c := range p {
		b = append(b, byte(c.Ident>>8), byte(c.Ident), byte(c.Site))
	}
	return b
}

// NewPos returns a position from the bytes. It doesn't validate the byte slice, so only
// pass into it valid bytes.
func NewPos(b []byte) []Identifier {
	p := []Identifier{}
	for i := 0; i < int(b[0]); i++ {
		offset := i*3 + 1
		ident := uint16(b[offset])<<8 + uint16(b[offset+1])
		site := uint8(b[offset+2])
		p = append(p, Identifier{ident, site})
	}
	return p
}