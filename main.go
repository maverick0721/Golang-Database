package main

import (
	"encoding/json"
	"fmt"
)

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

func main() {
	dir := "./"

	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	friends = []User{
		{"Yash", "21", "0987654321", "Google", Address{"Bangalore", "Karnataka", "India", "560001"}},
		{"Siddharth", "22", "0987654444", "Microsoft", Address{"Bangalore", "Karnataka", "India", "560001"}},
		{"Manav", "20", "0987653333", "Adobe", Address{"Bangalore", "Karnataka", "India", "560001"}},
		{"Shailendra", "23", "0987652222", "Netflix", Address{"Bangalore", "Karnataka", "India", "560001"}},
		{"Aman", "21", "0987651111", "Apple", Address{"Gwalior", "Madhya Pradesh", "India", "474011"}},
		{"Priyanshu", "21", "0987650000", "Amazon", Address{"Bangalore", "Karnataka", "India", "560001"}},
	}

	for _, value := range friends {
		db.Write("users", value.Name, User{
			Name:    value.Name,
			Age:     value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}

	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println(records)
}
