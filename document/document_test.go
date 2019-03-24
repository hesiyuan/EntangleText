package document

import (
	"fmt"
	"strings"
	"testing"
)

func TestDocNavigation(t *testing.T) {
	c := strings.Split("Entangle Text", "")
	clientID := uint8(1)
	mydoc := New(c, clientID) // static method
	p := End                  // constant
	i, _ := mydoc.Index(p)
	fmt.Println(i)
	p, _ = mydoc.Left(p)
	p, _ = mydoc.Right(p)
	p, _ = mydoc.Left(p)
	mydoc.DeleteLeft(p)
	fmt.Println(mydoc.Content())

}
