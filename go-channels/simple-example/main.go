package main

import (
	"fmt"
)

func main() {
	userch := make(chan string) // declaring a channel of type string

	// a channel is like a pipe where we can read/write data from/to.

	// userch <- "Bob"  // putting a value (string) to the channel
	// user := <-userch // reading from the channel
	// fmt.Println(user)

	/*
		1- A channel can be buffered or unbuffered. You can specify the size of buffer.
		2- A channel will block , if the buffer is full.
		3- In the above example (line 12 to 14) , since Bob is not in a routine, it is blocking .
		This gets solved once you put the channel read in a go routine (since go routines are non-blocking in nature)
	*/
	go func() {
		userch <- "Bob" // putting a value (string) to the channel
	}()
	user := <-userch // reading from the channel
	fmt.Println(user)

	// Now let's see a buffered channel
	userch2 := make(chan string, 2) // buffered channel of size 2

	userch2 <- "Alice"
	userch2 <- "Paul" // Unlike last example, this will work since now buffer is full so it is non-blocking
	user2 := <-userch2
	fmt.Println(user2) // Here one value will be read from the channel, it is consumed and 1 buffer space will become vacant
	user2 = <-userch2

	fmt.Println(user2)

	// sendMessage(userch2)
	// message := <-userch2
	// fmt.Println(message)

	chan3 := make(chan string, 3)
	sendMessage(chan3)
	readMessage(chan3)
}

func sendMessage(msgCh chan<- string) {
	// Here the param `chan<-` has specified that this method can only send to the channel msg
	// It cannot receive from the channel i.e. msg:=<- msgCh ; fmt.Println(msg) , this will throw error
	msgCh <- "Hello"
}

func readMessage(msgCh <-chan string) {
	//Here the param `<-chan` has specified that this method can only read to the channel msg
	msg := <-msgCh
	fmt.Println(msg + " Mayank")
}

/*
------ output -------

Bob
Alice
Paul
Hello Mayank

*/
