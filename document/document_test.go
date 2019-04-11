package document

import (
	"fmt"
	"strings"
	"testing"

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

func TestGenerationLongRP(t *testing.T) {
	clientID := uint8(68)
	lp := []Identifier{ // long lp case
		{0, 68},
	}

	rp := []Identifier{
		{0, 68},
		{0, 68},
		{2, 68},
	}

	for {

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
		assert.Equal(t, ComparePos(lp, p), int8(-1)) // p should be greater than lp
		assert.Equal(t, ComparePos(p, rp), int8(-1))

		rp = p
	}
}
