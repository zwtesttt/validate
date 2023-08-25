package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Validate struct {
	Lang          string
	MessageSep    string
	ErrMap        map[string]map[string]string
	isRequiredAll bool
}

type ValidationErrors struct {
	StructField string
	Field       string
	Value       interface{}
	Message     string
}

func New() *Validate {
	return &Validate{
		Lang:          "default",
		MessageSep:    "~",
		isRequiredAll: false,
		ErrMap:        make(map[string]map[string]string),
	}
}

func (v *Validate) ValidateStruct(stc interface{}) *ValidationErrors {
	validationErrors := ValidationErrors{}
	val := reflect.ValueOf(stc)
	//遍历结构体所有参数
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("validate")
		//判断规则是否为空
		if tag != "" {
			fieldValue := val.Field(i).Interface()
			rules := strings.Split(tag, ",")
			for _, rule := range rules {
				ruleParts := strings.Split(rule, "=")
				validatorName := ruleParts[0]
				validatorArgs := ruleParts[1:]
				v.customErrorMessage(field.Name, rule)
				if err := v.validateField(field.Name, validatorName, fieldValue, validatorArgs); err != nil {
					validationErrors.Message = err.Error()
					return &validationErrors
				}
			}
		}
	}
	return nil
}

func (v *Validate) customErrorMessage(field, option string) {
	ops := strings.Split(option, v.MessageSep)
	if len(ops) < 2 {
		return // 无效的选项格式
	}
	rules := strings.Split(ops[0], "=")
	if v.ErrMap[field] == nil {
		v.ErrMap[field] = make(map[string]string)
	}
	v.ErrMap[field][rules[0]] = ops[1:][0]

}

func (v *Validate) validateField(fieldName, validatorName string, value interface{}, args []string) error {

	cs := ValidationStruct{
		v.MessageSep, fieldName, validatorName, value, args, v.ErrMap,
	}
	if validationFunc, ok := validationFunctions[validatorName]; ok {
		if err := validationFunc(cs); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("No such rule:%s", validatorName)
	}

	field := reflect.ValueOf(value)
	if field.Kind() == reflect.Struct {
		if err := v.ValidateStruct(value); err != nil {
			return errors.New(err.Message)
		}
	}
	return nil
}
