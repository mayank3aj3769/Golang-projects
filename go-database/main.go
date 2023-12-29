package main

import (
	"fmt"
	"go-database/driver"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting gokuDB on :8080...")
	url := "http:localhost:8080/"
	_, err := http.Get(url)
	driver.DbStart()
	if err != nil {
		fmt.Println(err)
	}

	log.Println("========= Shutting down gokuDB ==========")

}
