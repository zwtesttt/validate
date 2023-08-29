package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	VARKEY = "one"
)

type Validate struct {
	Lang          string
	MessageSep    string
	ErrMap        map[string]map[string]string
	isRequiredAll bool
}

type ValidationErrors struct {
	StructField reflect.StructField
	Field       string
	Value       interface{}
	Message     string
}

func New() *Validate {
	return &Validate{
		Lang:          "en",
		MessageSep:    "~",
		isRequiredAll: false,
		ErrMap:        make(map[string]map[string]string),
	}
}

func (v *Validate) ValidateStruct(stc interface{}) error {
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
				ops := ruleParts[0]
				validatorName := strings.Split(ops, v.MessageSep)
				validatorArgs := ruleParts[1:]
				v.customErrorMessage(field.Name, rule)

				if err := v.validateField(field.Name, validatorName[0], fieldValue, validatorArgs); err != nil {
					//封装错误信息
					return &ValidationErrors{
						Message:     err.Error(),
						Value:       fieldValue,
						Field:       field.Name,
						StructField: field,
					}
				}
			}
		}
	}
	return nil
}

func (ve *ValidationErrors) Error() string {
	return ve.Message
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
		v.MessageSep, fieldName, validatorName, reflect.ValueOf(value), value, args, v.ErrMap,
	}
	if validationFunc, ok := validationFunctions[validatorName]; ok {
		if err := validationFunc(cs); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no such rule:%s", validatorName)
	}

	field := reflect.ValueOf(value)
	if field.Kind() == reflect.Struct {
		if err := v.ValidateStruct(value); err != nil {
			return errors.New(err.Error())
		}
	}
	return nil
}

// 添加自定义规则
func (v *Validate) AddValidationRule(validationName string, vsf ValidationFunc) error {
	if len(validationName) == 0 {
		return fmt.Errorf("null rule")
	}
	if _, ok := validationFunctions[validationName]; !ok {
		validationFunctions[validationName] = vsf
		return nil
	}

	return fmt.Errorf("existing rule")
}

// 校验单个变量
func (v *Validate) Var(arg interface{}, tag string) error {
	if len(tag) == 0 {
		return fmt.Errorf("null rule")
	}
	rules := strings.Split(tag, ",")
	for _, rule := range rules {
		ruleParts := strings.Split(rule, "=")
		ops := ruleParts[0]
		validatorName := strings.Split(ops, v.MessageSep)
		validatorArgs := ruleParts[1:]
		v.customErrorMessage(VARKEY, rule)
		if err := v.validateField(VARKEY, validatorName[0], arg, validatorArgs); err != nil {
			//封装错误信息
			return &ValidationErrors{
				Message: err.Error(),
				Value:   arg,
			}
		}

	}
	return nil
}
