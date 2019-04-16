package document

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestDocNavigation(t *testing.T) {
	c := strings.Split("Entangle Text", "")
	clientID := uint8(1)
	mydoc := NewDocument(c, clientID) // static method
	p := End                          // constant
	i, _ := mydoc.Index(p)
	i, _ = mydoc.Index(End)
	fmt.Println(i)
	p, _ = mydoc.Left(p)
	p, _ = mydoc.Right(p)
	p, _ = mydoc.Left(p)
	mydoc.DeleteLeft(p)
	fmt.Println(mydoc.Content())

}

func TestBatchTransfer(t *testing.T) {
	c := strings.Split("Entangle Text", "")
	clientID := uint8(1)
	doc1 := NewDocument(c, clientID) // static method
	//p := End                          // constant

	doc2 := Document{clientID: clientID}
	// doc2.insert(Start, "")
	// doc2.insert(End, "")
	// simualte batch transfer
	for _, e := range doc1.pairs {
		doc2.insert(e.pos, e.atom)
	}

	// check two doc are the same
	assert.Equal(t, len(doc1.pairs), len(doc2.pairs))

	for i, e := range doc1.pairs {
		if ComparePos(e.pos, doc2.pairs[i].pos) == 0 && e.atom == doc2.pairs[i].atom {
			continue
		} else {
			assert.Equal(t, 1, 0)
		}

	}
}

func TestConsistencyOne(t *testing.T) {
	c := strings.Split("abc", "")
	clientID := uint8(1)
	doc1 := NewDocument(c, clientID) // static method
	//p := End                          // constant
	doc2 := Document{clientID: clientID}
	// simualte batch transfer
	for _, e := range doc1.pairs {
		doc2.insert(e.pos, e.atom)
	}

	// peer 1 insert d between a and b.
	p1, flag := doc1.InsertLeft(doc1.pairs[2].pos, "d")
	// peer 2 delete a
	p2 := doc2.pairs[1].pos
	flag2 := doc2.delete(p2)

	if flag && flag2 == false {
		assert.Equal(t, 1, 0)
	}

	// transmiting pos are received
	doc2.insert(p1, "d")
	doc1.delete(p2)

	// check two doc are the same
	assert.Equal(t, len(doc1.pairs), len(doc2.pairs))

	for i, e := range doc1.pairs {
		if ComparePos(e.pos, doc2.pairs[i].pos) == 0 && e.atom == doc2.pairs[i].atom {
			continue
		} else {
			assert.Equal(t, 1, 0)
		}

	}
}

func TestRepeatedDeletes(t *testing.T) {
	c := strings.Split("abc", "")
	clientID := uint8(1)
	doc1 := NewDocument(c, clientID) // static method
	//p := End                          // constant
	doc2 := Document{clientID: clientID}
	// simualte batch transfer
	for _, e := range doc1.pairs {
		doc2.insert(e.pos, e.atom)
	}

	// peer 1 and peer 2 are deleting the same thing
	doc1_size, doc2_size := len(doc1.pairs), len(doc2.pairs)
	p1 := doc1.pairs[1].pos
	// peer 2 delete a
	p2 := doc2.pairs[1].pos
	flag := doc1.delete(p1)
	flag2 := doc2.delete(p2)

	if flag && flag2 == false {
		assert.Equal(t, 1, 0)
	}

	// transmiting pos are received
	doc2.delete(p1)
	doc1.delete(p2)

	// check two doc are the same and the length is decreased only by one
	assert.Equal(t, len(doc1.pairs), doc1_size-1)
	assert.Equal(t, len(doc2.pairs), doc2_size-1)

	for i, e := range doc1.pairs {
		if ComparePos(e.pos, doc2.pairs[i].pos) == 0 && e.atom == doc2.pairs[i].atom {
			continue
		} else {
			assert.Equal(t, 1, 0)
		}

	}
}

func TestHighConcurrency(t *testing.T) {
	c := strings.Split("abc", "")
	clientID := uint8(1)
	doc1 := NewDocument(c, clientID) // static method
	//p := End                          // constant
	doc2 := Document{clientID: clientID}
	// simualte batch transfer
	for _, e := range doc1.pairs {
		doc2.insert(e.pos, e.atom)
	}

	// peer 2 insert x between b and c
	p2, _ := doc2.InsertLeft(doc2.pairs[3].pos, "x")

	doc1_size := len(doc1.pairs)
	// GOAL: only testing race condition of doc1
	// peer 1 insert d between a and b while receiving a remote insert between b and c
	go func() { // simulate received RPC
		doc1.insert(p2, "x")
	}()

	p1, _ := doc1.InsertLeft(doc1.pairs[2].pos, "d")

	// transmiting pos are received
	doc2.delete(p1)
	doc1.delete(p2)

	time.Sleep(time.Second) // wait for routine to finish in a lazy way
	// check two doc are the same and the length is decreased only by one
	assert.Equal(t, len(doc1.pairs), doc1_size+2)

	fmt.Println(doc1.Content())
}

func TestGenerationPosLongLp(t *testing.T) {
	clientID := uint8(1)
	lp := []Identifier{ // long lp case
		{13627, 1},
		{65036, 1},
		{24224, 1},
	}

	rp := []Identifier{
		{13628, 1},
	}

	p, _ := GeneratePos(lp, rp, clientID)

	assert.Equal(t, ComparePos(lp, p), int8(-1)) // p should be greater than lp
	assert.Equal(t, ComparePos(p, rp), int8(-1))
}

func TestGenerationPos(t *testing.T) {
	clientID := uint8(68)
	lp := []Identifier{ // long lp case
		{65534, 68},
		{48896, 57},
		{65534, 68},
	}

	rp := []Identifier{
		{65534, 68},
		{48896, 68},
	}

	p, _ := GeneratePos(lp, rp, clientID) // p will be extended, though we don't have to

	assert.Equal(t, ComparePos(lp, p), int8(-1)) // p should be greater than lp
	assert.Equal(t, ComparePos(p, rp), int8(-1))
}

func TestGenerationPosLeftInsertion(t *testing.T) {
	clientID := uint8(68)
	lp := []Identifier{ // long lp case
		{56, 68},
		{31603, 68},
	}

	rp := []Identifier{
		{56, 68},
		{31603, 68},
		{1, 68},
	}

	p, _ := GeneratePos(lp, rp, clientID) // p will be extended, though we don't have to

	assert.Equal(t, ComparePos(lp, p), int8(-1)) // p should be greater than lp
	assert.Equal(t, ComparePos(p, rp), int8(-1))
}

func TestGenerationPosEqualLengthHasSpaces(t *testing.T) {
	clientID := uint8(68)
	lp := []Identifier{ // long lp case
		{56, 68},
		{31603, 68},
		{15, 68},
	}

	rp := []Identifier{
		{56, 68},
		{31603, 68},
		{278, 68},
	}

	p, _ := GeneratePos(lp, rp, clientID) // p will be extended, though we don't have to

	assert.Equal(t, ComparePos(lp, p), int8(-1)) // p should be greater than lp
	assert.Equal(t, ComparePos(p, rp), int8(-1))
}

func TestGenerationPosEqualLengthNoSpaces(t *testing.T) {
	clientID := uint8(68)
	lp := []Identifier{ // long lp case
		{56, 68},
		{31603, 68},
		{15, 68},
	}

	rp := []Identifier{
		{56, 68},
		{31603, 68},
		{16, 68},
	}

	p, _ := GeneratePos(lp, rp, clientID) // p will be extended, though we don't have to

	assert.Equal(t, ComparePos(lp, p), int8(-1)) // p should be greater than lp
	assert.Equal(t, ComparePos(p, rp), int8(-1))
}

func TestGenerationAnother(t *testing.T) {
	clientID := uint8(68)
	lp := []Identifier{ // long lp case
		{56, 68},
		{31603, 68},
		{65534, 68},
	}

	rp := []Identifier{
		{56, 68},
		{31603, 68},
		{65535, 68},
	}

	p, _ := GeneratePos(lp, rp, clientID) // p will be extended, though we don't have to

	assert.Equal(t, ComparePos(lp, p), int8(-1)) // p should be greater than lp
	assert.Equal(t, ComparePos(p, rp), int8(-1))
}

func TestGenerationSpecial(t *testing.T) {
	clientID := uint8(68)
	lp := []Identifier{ // long lp case
		{6623, 68},
		{65534, 68},
	}

	rp := []Identifier{
		{6624, 68},
		{62098, 68},
	}

	p, _ := GeneratePos(lp, rp, clientID) // p will be extended, though we don't have to

	assert.Equal(t, ComparePos(lp, p), int8(-1)) // p should be greater than lp
	assert.Equal(t, ComparePos(p, rp), int8(-1))
}

func TestGenerationSpecial2(t *testing.T) {
	clientID := uint8(68)
	lp := []Identifier{ // long lp case
		{6623, 68},
		{65534, 68},
	}

	rp := []Identifier{
		{6623, 68},
		{65535, 68},
	}

	p, _ := GeneratePos(lp, rp, clientID) // p will be extended, though we don't have to

	assert.Equal(t, ComparePos(lp, p), int8(-1)) // p should be greater than lp
	assert.Equal(t, ComparePos(p, rp), int8(-1))
}

func TestGenerationInsertLeft(t *testing.T) {
	clientID := uint8(68)
	lp := []Identifier{ // long lp case
		{0, 68},
	}

	rp := []Identifier{
		{0, 68},
		{0, 68},
		{2, 68},
	}

	counter := 0
	for {
		if counter == 100 {
			break
		}
		p, _ := GeneratePos(lp, rp, clientID) // p will be extended, though we don't have to

		fmt.Println("p: ")
		for _, e := range p {

			fmt.Printf("{")
			fmt.Print(e.Ident)
			fmt.Print(", ")
			fmt.Print(e.Site)
			fmt.Printf("},")
			fmt.Printf("\n")
		}

		fmt.Println("lp: ")
		for _, e := range lp {

			fmt.Printf("{")
			fmt.Print(e.Ident)
			fmt.Print(", ")
			fmt.Print(e.Site)
			fmt.Printf("},")
			fmt.Printf("\n")
		}

		fmt.Println("rp: ")
		for _, e := range rp {

			fmt.Printf("{")
			fmt.Print(e.Ident)
			fmt.Print(", ")
			fmt.Print(e.Site)
			fmt.Printf("},")
			fmt.Printf("\n")
		}

		counter = counter + 1
		assert.Equal(t, ComparePos(lp, p), int8(-1)) // p should be greater than lp
		assert.Equal(t, ComparePos(p, rp), int8(-1))

		rp = p
	}
}
