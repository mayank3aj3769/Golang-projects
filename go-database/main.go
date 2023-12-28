package main

import (
	"bufio"
	"fmt"
	"go-database/driver"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	db, err := driver.New("./", nil)
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1 - Add a new collection")
		fmt.Println("2 - Update values of an existing collection")
		fmt.Println("3 - Delete a collection")
		fmt.Println("4 - Show the list of existing collections")
		fmt.Println("0 - Exit")

		fmt.Print("Enter your choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			driver.AddNewCollection(db, reader)
		case "2":
			driver.UpdateCollection(db, reader)
		case "3":
			driver.DeleteCollection(db, reader)
		case "4":
			driver.ListCollections(db)
		case "0":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}
