package main

import (
	"fmt"
	"time"
)

type Server struct {
	users  map[string]string
	userch chan string
	quitch chan struct{}
}

func main() {
	server := NewServer()
	server.Start()

	// user, ok := <-server.userch // This can be used to see
	//	whether the channel has been closed by a different process or not

	// The following routing is just made to quit the loop for now.
	go func() {
		time.Sleep(2 * time.Second)
		//server.quitch <- struct{}{} // does the same thing as close
		close(server.quitch)
	}()

	// This blocks forever until at least one of the select cases(channels) becomes valid.
	// A default case can be used to make it non-blocking
	select {}

}
func NewServer() *Server { /* Constructor for the struct Server*/
	return &Server{
		users:  make(map[string]string),
		userch: make(chan string),
		quitch: make(chan struct{}),
	}
}

func (s *Server) Start() {
	go s.loop()
}

func (s *Server) loop() {

free:
	for {
		/* Here we are waiting until a user is put in channel. This blocking operation solves the race
		condition
		*/
		// user := <-s.userch
		// s.users[user] = user
		// fmt.Printf("Adding new user %s \n", user)
		select {
		case msg := <-s.userch:
			fmt.Println(msg)
		case <-s.quitch:
			fmt.Println("Server need to quit")
			break free
		default:
		}
	}
}

func (s *Server) addUser(user string) {
	s.users[user] = user
}

// func main() {

// 	ch1 := make(chan string, 2)
// 	ch1 <- "Hello"
// 	sendMessage(ch1)
// 	receiveMessage(ch1)

// 	sr := NewServer()
// 	//ch2 := make(chan string, 2)
// 	sr.userch <- "Hey"
// 	sr.addUser("Hey ")
// 	sr.users["Hello"] = "Mayank"
// 	user := <-sr.userch
// 	fmt.Println(user)
// 	sr.users["Hi"] = "Bob"

// 	for i := 0; i > len(sr.users); i++ {
// 		fmt.Println(sr.users)
// 	}
// }

func sendMessage(msgCh chan<- string) {
	msgCh <- " Mayank"
}
func receiveMessage(msgCh <-chan string) {
	msg := <-msgCh
	msg += <-msgCh
	fmt.Println(msg)
}
