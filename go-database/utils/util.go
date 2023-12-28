package utils

import (
	"encoding/json"
)

type Address struct {
	City    string
	Pincode json.Number
	State   string
	Country string
}

type User struct {
	Name    string
	Age     json.Number
	Company string
	Address Address
}
