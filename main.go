package main

import (
	"fmt"
	"validate/validator"
)

type User struct {
	Name    string  `validate:"required~不能为空"`
	Age     int64   `validate:"required~不能为空,min=18~最小值为%d"`
	Address Address `validate:"required~不能为空"`
	Email   string  `validate:"email"`
}

type Address struct {
	Street  string `validate:"required"`
	City    string `validate:"required~不能为空"`
	Country string `validate:"required~不能为空,min=666"`
}

func main() {
	address := Address{
		Street:  "ffff",
		City:    "6666",
		Country: "ssssss",
	}
	user := User{Name: "ssss", Age: 17, Address: address, Email: "123sss@qq.com"}
	validate := validator.New()
	var name string = ""

	if err := validate.Var(name, "required~36333"); err != nil {
		if _, ok := err.(*validator.ValidationErrors); ok {
			fmt.Println("Validation failed:", err.Error())
			return
		}

	}
	if err := validate.ValidateStruct(user); err != nil {
		if _, ok := err.(*validator.ValidationErrors); ok {
			fmt.Println("Validation failed:", err.Error())
			return
		}

	}
	fmt.Println("Validation passed.")

}

func validateRequired2(vs validator.ValidationStruct) error {
	//// 验证逻辑
	if vs.Value == "" || vs.Value == 0 || vs.Value == nil {
		errMsg := vs.MsgMap[vs.FieldName][vs.ValidatorName]
		if errMsg == "" {
			errMsg = "The value must66666"
		}
		return fmt.Errorf(errMsg)
	}
	return nil
}
