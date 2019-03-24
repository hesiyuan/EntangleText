package main

// all general error handling code goes here

import (
	"fmt"
	"os"
)

// If error is non-nil, print it out and halt.
func checkError(err error) {
	if err != nil {
		//fmt.Fprintf(os.Stderr, "MyError: ", err.Error())
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}
}
