package driver

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func DbStart() {

	reader := bufio.NewReader(os.Stdin)
	db, err := New("./", nil)
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1 - Add a new collection")
		fmt.Println("2 - Update values of an existing collection")
		fmt.Println("3 - Delete a collection")
		fmt.Println("4 - View a collection ")
		fmt.Println("5 - Show the list of existing collections")
		fmt.Println("0 - Exit")

		fmt.Print("Enter your choice: \n")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			AddNewCollection(db, reader)
			fmt.Println("---------")
		case "2":
			UpdateCollection(db, reader)
			fmt.Println("---------")
		case "3":
			DeleteCollection(db, reader)
			fmt.Println("---------")
		case "4":
			ViewCollection(db, reader)
			fmt.Println("---------")
		case "5":
			ListCollections(db)
			fmt.Println("---------")
		case "0":
			fmt.Println("Exiting...")

			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}
