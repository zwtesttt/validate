package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidationFunc func(ValidationStruct) error

type ValidationStruct struct {
	messageSep    string
	FieldName     string
	ValidatorName string
	Field         reflect.Value
	Value         interface{}
	Args          []string
	MsgMap        map[string]map[string]string
}

var validationFunctions = map[string]ValidationFunc{
	"required": validateRequired,
	"email":    validateEmail,
	"min":      validateMin,
	// 添加其他验证函数
}

func validateRequired(vs ValidationStruct) error {
	// 验证逻辑
	if vs.Value == "" || vs.Value == 0 || vs.Value == nil {
		errMsg := vs.MsgMap[vs.FieldName][vs.ValidatorName]
		if errMsg == "" {
			errMsg = "The value must"
		}
		return fmt.Errorf(errMsg)
	}
	return nil
}

func validateEmail(vs ValidationStruct) error {
	if !emailRegex.MatchString(fmt.Sprintf("%s", vs.Value)) {
		errMsg := vs.MsgMap[vs.FieldName][vs.ValidatorName]
		if errMsg == "" {
			errMsg = "email format error"
		}
		return fmt.Errorf(errMsg)
	}
	return nil
}

func validateMin(vs ValidationStruct) error {
	// 验证逻辑
	if len(vs.Args) == 0 {
		return fmt.Errorf("min validation requires an argument")
	}
	ops := strings.Split(vs.Args[0], vs.messageSep)
	min, err := strconv.Atoi(ops[0])
	if err != nil {
		return fmt.Errorf("%s Parameter error", vs.ValidatorName)
	}
	if intValue, ok := vs.Value.(int); ok {
		if intValue < min {
			errMsg := vs.MsgMap[vs.FieldName][vs.ValidatorName]
			if !strings.Contains(errMsg, "%d") {
				return fmt.Errorf(errMsg)
			}

			if errMsg == "" {
				errMsg = "value min is %d"
			}
			return fmt.Errorf(errMsg, min)
		}
	}
	return nil
}
