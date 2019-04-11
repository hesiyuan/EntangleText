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
