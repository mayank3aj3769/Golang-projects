package main

import (
	"encoding/json"
	"fmt"
	"go-database/driver"
	"go-database/utils"
)

const Version = "1.0.0" //version of db

func main() {
	dir := "./"
	db, err := driver.New(dir, nil)

	if err != nil {
		fmt.Println("Error ", err)
	}
	employees := []utils.User{
		{Name: "Mayank", Age: "27", Company: "Google", Address: utils.Address{City: "SF", Pincode: "30020",
			State: "CA", Country: "USA"}},
		{Name: "Raj", Age: "27", Company: "Google", Address: utils.Address{City: "SF", Pincode: "30021",
			State: "LA", Country: "USA"}},
		{Name: "John", Age: "24", Company: "Meta", Address: utils.Address{City: "NYC", Pincode: "30022",
			State: "NY", Country: "USA"}},
		{Name: "Doe", Age: "27", Company: "Microsoft", Address: utils.Address{City: "CH", Pincode: "30023",
			State: "IL", Country: "USA"}},
		{Name: "Paul", Age: "21", Company: "Netflix", Address: utils.Address{City: "MI", Pincode: "30024",
			State: "FL", Country: "USA"}},
		{Name: "Mike", Age: "29", Company: "Apple", Address: utils.Address{City: "AS", Pincode: "30025",
			State: "TX", Country: "USA"}},
	}

	for _, value := range employees {
		db.Write("users", value.Name, utils.User{
			Name:    value.Name,
			Age:     value.Age,
			Company: value.Company,
			Address: value.Address,
		})
	}
	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error ", err)
	}
	fmt.Println(records)

	allrecords := []utils.User{}

	for _, f := range records {
		employee := utils.User{}
		if err := json.Unmarshal([]byte(f), &employee); err != nil {
			fmt.Println("error ", err)
		}
		allrecords = append(allrecords, employee)
	}
	fmt.Println(allrecords) // fetch all records , add to a slice and print

	// if err := db.Delete("users", "John"); err != nil {
	// 	fmt.Println("Error ", err)
	// }

	// if err := db.Delete("users", ""); err != nil {
	// 	fmt.Println("Error ", err)
	// }
}
