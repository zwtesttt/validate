package main

import (
	"fmt"
	"validate/validator"
)

type User struct {
	Name    string  `validate:"required"`
	Age     int     `validate:"min=18~最小值为%d"`
	Address Address `validate:"required"`
}

type Address struct {
	Street  string `validate:"required"`
	City    string `validate:"required"`
	Country string `validate:"required~不能为空"`
}

func main() {
	address := Address{
		Street:  "ff",
		City:    "6666",
		Country: "",
	}
	user := User{Name: "ssss", Age: 18, Address: address}
	validate := validator.New()
	validate.MessageSep = "~"

	if err := validate.ValidateStruct(user); err != nil {
		fmt.Println("Validation failed:", err.Message)
	} else {
		fmt.Println("Validation passed.")
	}
}
