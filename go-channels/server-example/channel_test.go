package main

import (
	"fmt"
	"testing"
)

// go test ./... -v --race
func TestAddUser(t *testing.T) {
	server := NewServer()
	server.Start()
	/*
		This leads to race condition problem
	*/
	for i := 0; i < 10; i++ {
		go func(i int) {
			// server.addUser(fmt.Sprintf("user_%d", i)) --> Gives race condition
			server.userch <- fmt.Sprintf("user_%d", i)
		}(i)
	}
	fmt.Println("The loop is done")
}
